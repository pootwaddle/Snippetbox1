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

func (db *Database) LatestSnippets() (Snippets, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	// Use the QueryRow() method on the embedded connection pool to execute our
	// SQL statement. This results a sql.Rows resultset containing the result of
	// our query.
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	// IMPORTANTLY we defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before LatestSnippets() returns. Closing a
	// resultset is really important. As long as a resultset is open it will
	// keep the underlying database connection open. So if something goes wrong
	// in this method and the resultset isn't closed, it can rapidly lead to all
	// the connections in your pool being used up. Another gotcha is that the
	// defer statement should come *after* you check for an error from
	// db.Query(). Otherwise, if db.Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty Snippets object (remember that this is just a slice of
	// the type []*Snippet).
	snippets := Snippets{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all of the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet object.
		s := &Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice.
	return snippets, nil
}

func (db *Database) InsertSnippet(title, content, expires string) (int, error) {
	// Write the SQL statement we want to execute.
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? SECOND))`

	// Use the db.Exec() method to execute the statement snippet, passing in values
	// for our (untrusted) title, content and expiry placeholder parameters in
	// exactly the same way that we did with the QueryRow() method. This returns
	// a sql.Result object, which contains some basic information about what
	// happened when the statement was executed.
	result, err := db.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result object to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned is of type int64, so we convert it to an int type for
	// returning from our Insert function.
	return int(id), nil
}

func ExampleTransaction(db *sql.DB) error {
	// Calling the Begin() method on the connection pool creates a new sql.Tx
	// object, which represents the in-progress database transaction.
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Call Exec() on the transaction, passing in your statement and any
	// parameters. It's important to notice that tx.Exec() is called on the
	// transaction object we just created, NOT the connection pool. Although we're
	// using tx.Exec() here you can also use tx.Query() and tx.QueryRow() in
	// exactly the same way.
	_, err = tx.Exec("INSERT INTO ...")
	if err != nil {
		// If there is any error, we call the tx.Rollback() method on the
		// transaction. This will abort the transaction and no changes will be
		// made to the database.
		tx.Rollback()
		return err
	}

	// Carry out another transaction in exactly the same way.
	_, err = tx.Exec("UPDATE ...")
	if err != nil {
		tx.Rollback()
		return err
	}

	// If there are no errors, the statements in transaction can be commited
	// to the database with the tx.Commit() method. It's really important to ALWAYS
	// call either Rollback() or Commit() before your function returns. If you
	// don't the connection will stay open and not be returned to the connection
	// pool. This can lead to hitting you maximum connection limit/running out of
	// resources.
	err = tx.Commit()
	return err
}
