package dao

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
)

const (
	insertUserSQL        = "INSERT INTO user(id, email, name, location, password) VALUES($1, $2, $3, $4, $5)"
	selectUserSQL        = "SELECT id, email, name, location FROM user WHERE email=? AND password=?"
	selectUserByEmailSQL = "SELECT id, email, name, location FROM user WHERE email=?"
)

var (
	errEmptyUser = errors.New("empty user not allowed")

	// ErrUserNotFound is a flag error to indicate a not found user error
	ErrUserNotFound = errors.New("user not found")
)

// User represents a user record
type User struct {
	ID           int
	Email        string
	Name         string
	Location     string
	PasswordHash string `json:"-"`
}

// InsertUser creates a new user
func (d *DAO) InsertUser(ctx context.Context, u *User) error {
	if u == nil {
		return errEmptyUser
	}
	result, err := d.db.ExecContext(ctx, insertUserSQL, u.ID, u.Email, u.Name, u.Location, u.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	countRows, _ := result.RowsAffected()
	if countRows != 1 {
		return fmt.Errorf("unexpected error inserting user: %+v", *u)
	}
	return nil
}

// GetUserByEmailPasswd fetchs a user by email and password
func (d *DAO) GetUserByEmailPasswd(ctx context.Context, email, textPasswd string) (*User, error) {
	h := sha256.New()
	h.Write([]byte(textPasswd))
	hashPasswd := fmt.Sprintf("%x", h.Sum(nil))

	u := new(User)
	err := d.db.QueryRowContext(ctx, selectUserSQL, email, hashPasswd).Scan(&u.ID, &u.Email, &u.Name, &u.Location)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrUserNotFound
	case err != nil:
		return nil, fmt.Errorf("select user error: %w", err)
	}
	return u, nil
}

// GetUserByEmail fetchs a user by email
func (d *DAO) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := new(User)
	err := d.db.QueryRowContext(ctx, selectUserByEmailSQL, email).Scan(&u.ID, &u.Email, &u.Name, &u.Location)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrUserNotFound
	case err != nil:
		return nil, fmt.Errorf("select user error: %w", err)
	}
	return u, nil
}
