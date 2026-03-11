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
	return 0, nil
}

// Returns snippets based on ID
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// Returns 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
