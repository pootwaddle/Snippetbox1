package models

import (
	"time"
)

// Declare a Database type (for now it's just an empty struct).
type Database struct{}

// Implement a GetSnippet() method on the Database type. For now, this just returns
// some dummy data, but later we'll update it to query our MySQL database for a
// snippet with a specific ID. In particular, it returns a dummy snippet if the id
// passed to the method equals 123, or returns nil otherwise.
func (db *Database) GetSnippet(id int) (*Snippet, error) {
	if id == 123 {
		snippet := &Snippet{
			ID:      id,
			Title:   "Example title",
			Content: "Example content",
			Created: time.Now(),
			Expires: time.Now(),
		}
		return snippet, nil
	}

	return nil, nil
}
