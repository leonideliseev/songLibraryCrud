package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/songLibraryCrud/internal/sqlc/queries"
	"github.com/leonideliseev/songLibraryCrud/internal/utils/convert/song"
	"github.com/leonideliseev/songLibraryCrud/models"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
)

type SongsPostgres struct {
	log *logging.Logger
	builder squirrel.StatementBuilderType
	conn *pgxpool.Pool
	queries *queries.Queries
}

func NewSongsPostgres(conn *pgxpool.Pool, log *logging.Logger) *SongsPostgres {
    queries := queries.New(conn)

	return &SongsPostgres{
		log: log,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		conn: conn,
		queries: queries,
	}
}

// для этой функции используется squirrel, а не sqlc, так как мне нужна динамическая генерация sql запроса
func (d *SongsPostgres) GetAll(ctx context.Context, limit, offset int, pagModel *models.Song) ([]models.Song, error) {
	query := d.builder.Select("*").From("songs").Limit(uint64(limit)).Offset(uint64(offset))
	query = addWhereWithCondition(query, "group_name", pagModel.GroupName)
	query = addWhereWithCondition(query, "name", pagModel.Name)
	query = addWhereWithCondition(query, "release_date", pagModel.ReleaseDate)
	query = addWhereWithCondition(query, "text", pagModel.Text)
	query = addWhereWithCondition(query, "link", pagModel.Link)

	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := d.conn.Query(ctx, q, args...)
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

func (d *SongsPostgres) CreateSong(ctx context.Context, s *models.Song) (*models.Song, error) {
	createSong := songConvert.FromAppToQuery(s)

	songQuery, err := d.queries.CreateSong(ctx, queries.CreateSongParams{
		GroupName: createSong.GroupName,
		Name: createSong.Name,
		ReleaseDate: createSong.ReleaseDate,
		Text: createSong.Text,
		Link: createSong.Link,
	})
	if err != nil {
		return nil, err
	}

	song := songConvert.FromQueryToApp(&songQuery)

	return song, nil
}

func (d *SongsPostgres) GetSong(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	uuid := toUUID(id)

	songQuery, err := d.queries.GetSong(ctx, uuid)
	if err != nil {
		return nil, err
	}

	song := songConvert.FromQueryToApp(&songQuery)

	return song, nil
}

func (d *SongsPostgres) UpdateSong(ctx context.Context, updatedData *models.Song) (*models.Song, error) {
	updateSong := songConvert.FromAppToQuery(updatedData)

	songQuery, err := d.queries.UpdateSong(ctx, queries.UpdateSongParams{
		ID: updateSong.ID,
		GroupName: updateSong.GroupName,
		Name: updateSong.Name,
		ReleaseDate: updateSong.ReleaseDate,
		Text: updateSong.Text,
		Link: updateSong.Link,
	})
	if err != nil {
		return nil, err
	}

	song := songConvert.FromQueryToApp(&songQuery)

	return song, nil
}

func (d *SongsPostgres) DeleteSong(ctx context.Context, id uuid.UUID) error {
	uuid := toUUID(id)

	return d.queries.DeleteSong(ctx, uuid)
}

func toUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}
