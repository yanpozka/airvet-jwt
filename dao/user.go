package dao

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
)

const (
	insertUserSQL = "insert into user(id, email, name, location, password) values($1, $2, $3, $4, $5)"
	selectUserSQL = "SELECT id, email, name, location FROM user WHERE email=? AND password=?"
)

var (
	emptyUserErr = errors.New("empty user not allowed")

	UserNotFoundErr = errors.New("user not found")
)

type User struct {
	ID           int
	Email        string
	Name         string
	Location     string
	PasswordHash string
}

func (d *DAO) InsertUser(ctx context.Context, u *User) error {
	if u == nil {
		return emptyUserErr
	}
	result, err := d.db.ExecContext(ctx, insertUserSQL, u.ID, u.Email, u.Name, u.Location, u.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to insert user: %+w", err)
	}

	countRows, _ := result.RowsAffected()
	if countRows != 1 {
		return fmt.Errorf("unexpected error inserting user: %+v", *u)
	}
	return nil
}

func (d *DAO) GetUserByEmailPasswd(ctx context.Context, email, textPasswd string) (*User, error) {
	h := sha256.New()
	h.Write([]byte(textPasswd))
	hashPasswd := fmt.Sprintf("%x", h.Sum(nil))
	u := new(User)
	err := d.db.QueryRowContext(ctx, selectUserSQL, email, hashPasswd).Scan(&u.ID, &u.Email, &u.Name, &u.Location)
	switch {
	case err == sql.ErrNoRows:
		return nil, UserNotFoundErr
	case err != nil:
		return nil, fmt.Errorf("select user error: %w", err)
	}
	return u, nil
}
