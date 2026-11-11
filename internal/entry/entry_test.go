package entry

import (
	"database/sql"
	"slices"
	"testing"
	"time"
)

func TestFields(t *testing.T) {
	time.Local = time.UTC // make Local() deterministic on every OS

	created := time.Date(2026, 10, 17, 14, 5, 0, 0, time.UTC)
	created.Format(timeLayout)

	tests := []struct {
		name string
		in   Entry
		want []string
	}{
		{
			name: "tag and project set",
			in: Entry{
				Id:        7,
				Message:   "hello",
				Tag:       sql.NullString{String: "work", Valid: true},
				Project:   sql.NullString{String: "devlog", Valid: true},
				CreatedAt: created,
			},
			want: []string{"7", "2026-10-17 14:05", "work", "devlog", "hello"},
		},
		{
			name: "tag and project null",
			in:   Entry{Id: 1, Message: "plain", CreatedAt: created},
			want: []string{"1", "2026-10-17 14:05", "-", "-", "plain"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Fields()
			if !slices.Equal(got, tt.want) {
				t.Errorf("Fields() = %v, want %v", got, tt.want)
			}
		})
	}
}
