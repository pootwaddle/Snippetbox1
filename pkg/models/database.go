package models

import (
	"database/sql"
)

// Declare a Database type (for now it's just an empty struct).
type Database struct {
	*sql.DB
}

/*
DB.Query() is used for SELECT queries which return multiple rows.
DB.QueryRow() is used for SELECT queries which return a single row.
DB.Exec() is used for statements which don't return rows (like INSERT and DELETE).
*/

// Implement a GetSnippet() method on the Database type. For now, this just returns
// some dummy data, but later we'll update it to query our MySQL database for a
// snippet with a specific ID. In particular, it returns a dummy snippet if the id
// passed to the method equals 123, or returns nil otherwise.
func (db *Database) GetSnippet(id int) (*Snippet, error) {
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is why it's surrounded with backticks instead
	// of normal double quotes).
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Use the QueryRow() method on the embedded connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result returned by the database.
	row := db.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &Snippet{}

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement. If our query returned no rows, then
	// row.Scan() will return a sql.ErrNoRows error. We check for that and return
	// nil instead of a Snippet object.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippet object.
	return s, nil
}
