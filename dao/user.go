package dao

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
)

const (
	insertUserSQL = "INSERT INTO user(id, email, name, location, password, privatekey, publickey) VALUES($1, $2, $3, $4, $5, $6, $7)"
	selectUserSQL = "SELECT id, email, name, location, privatekey, publickey FROM user WHERE email=? AND password=?"
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
	PasswordHash string
	PrivateKey   string
	PublicKey    string

	privateRSAKey *rsa.PrivateKey
}

func (u *User) GetRSAKey() (*rsa.PrivateKey, error) {
	if u.privateRSAKey == nil {
		// parse the private key once
		block, _ := pem.Decode([]byte(u.PrivateKey))
		if block == nil || block.Type != rsaPrivateKeyType {
			log.Fatal("failed to decode PEM block containing private key")
		}
		pkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		u.privateRSAKey = pkey
	}
	return u.privateRSAKey, nil
}

// InsertUser creates a new user
func (d *DAO) InsertUser(ctx context.Context, u *User) error {
	if u == nil {
		return errEmptyUser
	}
	result, err := d.db.ExecContext(ctx, insertUserSQL, u.ID, u.Email, u.Name, u.Location, u.PasswordHash, u.PrivateKey, u.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to insert user: %+w", err)
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
	err := d.db.QueryRowContext(ctx, selectUserSQL, email, hashPasswd).Scan(&u.ID, &u.Email, &u.Name, &u.Location, &u.PrivateKey, &u.PublicKey)
	switch {
	case err == sql.ErrNoRows:
		return nil, ErrUserNotFound
	case err != nil:
		return nil, fmt.Errorf("select user error: %w", err)
	}
	return u, nil
}
