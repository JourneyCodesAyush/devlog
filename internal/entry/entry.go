package entry

import (
	"database/sql"
	"fmt"
	"time"
)

type Entry struct {
	Id        int
	Message   string
	Tag       sql.NullString
	Project   sql.NullString
	CreatedAt time.Time
}

func (e *Entry) String() string {
	tag := ""
	if e.Tag.Valid {
		tag = e.Tag.String
	}
	return fmt.Sprintf("%s\t%s\t%s", e.CreatedAt.Local(), e.Message, tag)
}
