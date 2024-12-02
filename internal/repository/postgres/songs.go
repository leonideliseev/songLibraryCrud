package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/leonideliseev/songLibraryCrud/internal/repository/repoerr"
	"github.com/leonideliseev/songLibraryCrud/models"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/leonideliseev/songLibraryCrud/pkg/postgresql"
)

type SongsPostgres struct {
	log *logging.Logger
	builder squirrel.StatementBuilderType
	conn postgresql.Conn
}

func NewSongsPostgres(conn postgresql.Conn, log *logging.Logger) *SongsPostgres {
	defer log.Info("repository implementation inited successfully")
	return &SongsPostgres{
		log: log,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		conn: conn,
	}
}

func (r *SongsPostgres) GetAll(ctx context.Context, limit, offset int, pagModel *models.Song) ([]models.Song, error) {
	query := selectALL(r.builder).
		From(songsTable).
		Limit(uint64(limit)).
		Offset(uint64(offset))
	query = addWhereWithCondition(query, group_name_F, pagModel.GroupName)
	query = addWhereWithCondition(query, name_F, pagModel.Name)
	query = addWhereWithDateCondition(query, release_date_F, pagModel.ReleaseDate)
	query = addWhereWithCondition(query, text_F, pagModel.Text)
	query = addWhereWithCondition(query, link_F, pagModel.Link)

	q, args, err := query.ToSql()
	if err != nil {
		r.log.WithError(err).Error("failed to build query")
		return nil, err
	}

	r.log.WithField("query", q).WithField("args", args).Debug("get all songs query")

	rows, err := r.conn.Query(ctx, q, args...)
	if err != nil {
		r.log.WithError(err).
			WithField("limit", limit).
			WithField("offset", offset).
			Error("failed to get all songs")
		return nil, err
	}
	defer rows.Close()

	songs, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Song])
	if err != nil {
		r.log.WithError(err).Error("failed to collect rows")
		return nil, err
	}

	return songs, nil
}

func (r *SongsPostgres) Create(ctx context.Context, s *models.Song) (*models.Song, error) {
	q, args, err := r.builder.
		Insert(songsTable).
		Columns(group_name_F, name_F, release_date_F, text_F, link_F).
		Values(s.GroupName, s.Name, s.ReleaseDate, s.Text, s.Link).
        Suffix(fmt.Sprintf("RETURNING %s", ALL)).
		ToSql()
	if err != nil {
		r.log.WithError(err).Error("failed to build query")
		return nil, err
	}

	r.log.WithField("query", q).WithField("args", args).Debug("get create song query")

	row := r.conn.QueryRow(ctx, q, args...)

	var song models.Song
	if err := scanALL(row, &song); err != nil {
		var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == uniqError {
            r.log.WithError(err).
				WithField(name_F, s.Name).
				WithField(group_name_F, s.GroupName).
				Warn("song already exists")
            return nil, repoerr.ErrAlreadyExists
        }

		r.log.WithError(err).Error("failed to create song")
		return nil, err
	}

	r.log.WithField(id_F, s.ID).Info("song created successfully")
	return &song, nil
}

func (r *SongsPostgres) GetById(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	q, args, err := selectALL(r.builder).
		From(songsTable).
		Where(squirrel.Eq{id_F: id}).
		ToSql()
	if err != nil {
		r.log.WithError(err).Error("failed to build query")
		return nil, err
	}

	r.log.WithField("query", q).WithField("args", args).Debug("get song query")

	row := r.conn.QueryRow(ctx, q, args...)

	var song models.Song
	if err := scanALL(row, &song); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.WithError(err).WithField(id_F, id).Warn("song not found")
			return nil, repoerr.ErrNotFound
		}

		r.log.WithError(err).WithField(id_F, id).Error("failed to get song")
		return nil, err
	}

	return &song, nil
}

func (r *SongsPostgres) UpdateById(ctx context.Context, s *models.Song) (*models.Song, error) {
	q, args, err := r.builder.
		Update(songsTable).
		Set(group_name_F, s.GroupName).
		Set(name_F, s.Name).
		Set(release_date_F, s.ReleaseDate).
		Set(text_F, s.Text).
		Set(link_F, s.Link).
		Where(squirrel.Eq{id_F: s.ID}).
		ToSql()
	if err != nil {
		r.log.WithError(err).Error("failed to build query")
		return nil, err
	}

	r.log.WithField("query", q).WithField("args", args).Debug("get update song query")

	commandTag, err := r.conn.Exec(ctx, q, args...)
	if err != nil {
		var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == uniqError {
            r.log.WithError(err).
                WithField(group_name_F, s.GroupName).
                WithField(name_F, s.Name).
                Warn("song update conflict: unique constraint violated")
            return nil, repoerr.ErrAlreadyExists
        }

		r.log.WithError(err).WithField(id_F, s.ID).Error("failed to update song")
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		r.log.WithField(id_F, s.ID).Warn("no song found to update")
    	return nil, repoerr.ErrNotFound
	}

	r.log.WithField(id_F, s.ID).Info("song updated successfully")

	return s, nil
}

func (r *SongsPostgres) DeleteById(ctx context.Context, id uuid.UUID) error {
	q, args, err := r.builder.
		Delete(songsTable).
		Where(squirrel.Eq{id_F: id}).
		ToSql()
	if err != nil {
		r.log.WithError(err).Error("failed to build query")
		return err
	}

	r.log.WithField("query", q).WithField("args", args).Debug("get delete song query")

	commandTag, err := r.conn.Exec(ctx, q, args)
	if err != nil {
		r.log.WithError(err).WithField(id_F, id).Error("failed to delete song")
		return err
	}

	if commandTag.RowsAffected() == 0 {
		r.log.WithField(id_F, id).Warn("no song found to delete")
    	return repoerr.ErrNotFound
	}

	r.log.WithField(id_F, id).Info("song deleted successfully")
	return nil
}
