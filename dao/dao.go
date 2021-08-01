package dao

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	createUserTableAndEmptySQL = `CREATE TABLE IF NOT EXISTS user (
	id integer not null primary key,
	email text,
	name text,
	location text,
	password text); DELETE FROM user;`

	createJWKSAndEmptyTable = `CREATE TABLE IF NOT EXISTS jwks (
	privatekey text,
	publickey text,
	expiresAt integer); DELETE FROM jwks;`
)

// DAO is a Data Access Object abstracts DB manipulation
type DAO struct {
	db *sql.DB
}

// NewDAO creates a new DAO object
func NewDAO(path string) (*DAO, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	return &DAO{db: db}, nil
}

func (d *DAO) Close() {
	d.db.Close()
}

// InitDB creates and seeds the user table
func (d *DAO) InitDB(ctx context.Context) error {
	if _, err := d.db.ExecContext(ctx, createUserTableAndEmptySQL); err != nil {
		return fmt.Errorf("failed to create table user: %w", err)
	}
	if _, err := d.db.ExecContext(ctx, createJWKSAndEmptyTable); err != nil {
		return fmt.Errorf("failed to create table user: %w", err)
	}

	privatekey, publickey, err := generatePrivatePublicKeyPair()
	if err != nil {
		return err
	}
	if err := d.InsertJWKS(ctx, &JWK{PrivateKey: privatekey, PublicKey: publickey}); err != nil {
		return err
	}

	users := []*User{
		{
			ID:           1,
			Email:        "admin@airvet.com",
			Name:         "Admin",
			Location:     "somewhere",
			PasswordHash: "a9f4edc6c0f72ed3156a540dab48828f196066b32f9e41469b61069dcf62b80b", // "Admin-pass"
		},
		{
			ID:           2,
			Email:        "coolvet@airvet.com",
			Name:         "Cool Vet",
			Location:     "Best Pet Veterinary Clinic",
			PasswordHash: "0b04099717ab5a1bf87bccf2b1253bbf1206cde80c91a6cc30d62a3d5d82cae5", // "Cool_pass123"
		},
	}
	for _, u := range users {
		if err := d.InsertUser(ctx, u); err != nil {
			return err
		}
	}
	return nil
}
