package database

import "database/sql"

func NewStorage(url string) (*Queries, error) {
	conn, err := sql.Open("postgres", url)
	db := New(conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
