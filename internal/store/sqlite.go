package store

import (
	"database/sql"
	"fmt"

	"github.com/journeycodesayush/devlog/internal/entry"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS entries (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	message TEXT NOT NULL,
	tag TEXT,
	project TEXT,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_entries_created_at ON entries(created_at);
CREATE INDEX IF NOT EXISTS idx_entries_tag ON entries(tag);
`

func initializeDatabase(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err

}

func New(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	err = initializeDatabase(db)

	if err != nil {
		return nil, fmt.Errorf("initializing schema: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Add(message string, tag string) error {
	var tagVal sql.NullString
	if tag != "" {
		tagVal = sql.NullString{String: tag, Valid: true}
	}
	_, err := s.db.Exec("INSERT INTO entries (message, tag) VALUES (?, ?)", message, tagVal)
	if err != nil {
		return fmt.Errorf("adding entry: %w", err)
	}
	return nil
}

func (s *SQLiteStore) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM entries WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("deleting entry: %w", err)
	}
	return nil
}

func (s *SQLiteStore) List(limit int) ([]entry.Entry, error) {
	query := "SELECT * FROM entries ORDER BY created_at DESC, id DESC"
	var args []any

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return make([]entry.Entry, 0), fmt.Errorf("querying entry: %w", err)
	}

	defer rows.Close()

	entries := make([]entry.Entry, 0)

	for rows.Next() {
		var ent entry.Entry
		err := rows.Scan(&ent.Id, &ent.Message, &ent.Tag, &ent.Project, &ent.CreatedAt)

		if err != nil {
			return make([]entry.Entry, 0), fmt.Errorf("querying entry: %w", err)
		}

		entries = append(entries, ent)
	}

	if err = rows.Err(); err != nil {
		return make([]entry.Entry, 0), fmt.Errorf("querying entry: %w", err)
	}

	return entries, nil
}

func (s *SQLiteStore) Search(keyword string) ([]entry.Entry, error) {
	query := "SELECT * FROM entries WHERE message LIKE ? ORDER BY created_at DESC, id DESC"
	rows, err := s.db.Query(query, "%"+keyword+"%")
	if err != nil {
		return make([]entry.Entry, 0), fmt.Errorf("querying entry: %w", err)
	}

	defer rows.Close()

	entries := make([]entry.Entry, 0)
	for rows.Next() {
		var ent entry.Entry
		if err := rows.Scan(&ent.Id, &ent.Message, &ent.Tag, &ent.Project, &ent.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning entry: %w", err)
		}
		entries = append(entries, ent)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reading rows: %w", err)
	}
	return entries, nil
}
