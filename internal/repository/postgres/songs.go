package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/songLibraryCrud/models"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
)

const (
	songsTable = "songs"
)

// songs field
const (
	id_F = "id"
	group_name_F = "group_name"
	name_F = "name"
	release_date_F = "release_date"
	text_F = "text"
	link_F = "link"
)

type SongsPostgres struct {
	log *logging.Logger
	builder squirrel.StatementBuilderType
	conn *pgxpool.Pool
}

func NewSongsPostgres(conn *pgxpool.Pool, log *logging.Logger) *SongsPostgres {
	return &SongsPostgres{
		log: log,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		conn: conn,
	}
}

func (r *SongsPostgres) GetAll(ctx context.Context, limit, offset int, pagModel *models.Song) ([]models.Song, error) {
	query := r.builder.Select("*").From(songsTable).Limit(uint64(limit)).Offset(uint64(offset))
	query = addWhereWithCondition(query, group_name_F, pagModel.GroupName)
	query = addWhereWithCondition(query, name_F, pagModel.Name)
	query = addWhereWithCondition(query, release_date_F, pagModel.ReleaseDate)
	query = addWhereWithCondition(query, text_F, pagModel.Text)
	query = addWhereWithCondition(query, link_F, pagModel.Link)

	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Song])
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func addWhereWithCondition(query squirrel.SelectBuilder, field, param string) squirrel.SelectBuilder {
	if param != "" {
		return query.Where(squirrel.ILike{field: "%" + param + "%"})
	}
	return query
}

func (r *SongsPostgres) Create(ctx context.Context, s *models.Song) (*models.Song, error) {
	query, args, err := r.builder.Insert(songsTable).
	Columns(group_name_F, name_F, release_date_F, text_F, link_F).
	Values(s.GroupName, s.Name, s.ReleaseDate, s.Text, s.Link).ToSql()
	if err != nil {
		return nil, err
	}

	row := r.conn.QueryRow(ctx, query, args...)
	var song models.Song
	err = row.Scan(&song)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		return nil, err
	}

	return &song, nil
}

func (r *SongsPostgres) GetById(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	q, args, err := r.builder.Select("*").From(songsTable).Where(squirrel.Eq{id_F: id}).ToSql()
	if err != nil {
		return nil, err
	}

	row := r.conn.QueryRow(ctx, q, args...)

	var song models.Song
	err = row.Scan(&song)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

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
		return nil, err
	}

	_, err = r.conn.Exec(ctx, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		return nil, err
	}

	return s, err
}

func (r *SongsPostgres) DeleteById(ctx context.Context, id uuid.UUID) error {
	q, args, err := r.builder.Delete(songsTable).Where(squirrel.Eq{id_F: id}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(ctx, q, args)
	if err != nil {
		return err
	}
	
	return nil
}
