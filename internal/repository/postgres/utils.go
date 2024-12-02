package postgres

import (
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/leonideliseev/songLibraryCrud/models"
)

const (
	songsTable = "songs"
)

// songs field
const (
	id_F           = "id"
	group_name_F   = "group_name"
	name_F         = "name"
	release_date_F = "release_date"
	text_F         = "text"
	link_F         = "link"
)

var (
	ALL = fmt.Sprintf("%s, %s, %s, %s, %s, %s", id_F, group_name_F, name_F, release_date_F, text_F, link_F)
)

const (
	uniqError = "23505"
)

func selectALL(b squirrel.StatementBuilderType) squirrel.SelectBuilder {
	return b.Select(id_F, group_name_F, name_F, release_date_F, text_F, link_F)
}

func scanALL(row pgx.Row, song *models.Song) error {
	return row.Scan(
        &song.ID,
        &song.GroupName,
        &song.Name,
        &song.ReleaseDate,
        &song.Text,
        &song.Link,
	)
}

func addWhereWithCondition(query squirrel.SelectBuilder, field, param string) squirrel.SelectBuilder {
	if param != "" {
		return query.Where(squirrel.ILike{field: "%" + param + "%"})
	}
	return query
}

func addWhereWithDateCondition(query squirrel.SelectBuilder, field string, param time.Time) squirrel.SelectBuilder {
	if !param.IsZero() {
		return query.Where(squirrel.Eq{field: param})
	}
	return query
}
