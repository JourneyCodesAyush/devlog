package entry

import (
	"database/sql"
	"strconv"
	"time"
)

const timeLayout = "2006-01-02 15:04"

type Entry struct {
	Id        int
	Message   string
	Tag       sql.NullString
	Project   sql.NullString
	CreatedAt time.Time
}

func (e *Entry) Fields() []string {
	tag := "-"
	project := "-"
	if e.Tag.Valid {
		tag = e.Tag.String
	}
	if e.Project.Valid {
		project = e.Project.String
	}

	return []string{
		strconv.Itoa(e.Id),
		e.CreatedAt.Local().Format(timeLayout),
		tag,
		project,
		e.Message,
	}
}
