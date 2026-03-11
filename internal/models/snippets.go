package models

import (
	"database/sql"
	"time"
)

// Snippet represents a single text snippet with its metadata.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel wraps a database connection pool and provides methods
// for interacting with the snippets table.
type SnippetModel struct {
	DB *sql.DB
}

// Inserting new snippet into the database
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

// Returns snippets based on ID
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// Returns 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
