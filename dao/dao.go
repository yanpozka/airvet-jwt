package dao

import (
	"context"
	"database/sql"
	"fmt"
)

const createTableSQL = `CREATE TABLE IF NOT EXISTS user (id integer not null primary key, email text, name text, location text, password text)`

// Data Access Object abstracts DB manipulation
type DAO struct {
	db *sql.DB
}

func NewDAO(path string) (*DAO, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %+w", err)
	}

	return &DAO{db: db}, nil
}

func (d *DAO) InitDB(ctx context.Context) error {
	if _, err := d.db.ExecContext(ctx, createTableSQL); err != nil {
		return fmt.Errorf("failed to create table user: %+w", err)
	}

	users := []*User{
		{
			ID:           1,
			Email:        "admin@airvet.com",
			Name:         "Admin",
			Location:     "somewhere",
			PasswordHash: "a9f4edc6c0f72ed3156a540dab48828f196066b32f9e41469b61069dcf62b80b", // "Admin-pass"
		},
	}
	for _, u := range users {
		if err := d.InsertUser(ctx, u); err != nil {
			return err
		}
	}
	return nil
}
