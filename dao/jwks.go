package dao

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
)

const (
	insertJWKSQL  = "INSERT INTO jwks(privatekey, publickey) VALUES($1, $2)"
	selectJWKSSQL = "SELECT privatekey, publickey FROM jwks"
)

type JWK struct {
	PrivateKey string
	PublicKey  string

	privateRSAKey *rsa.PrivateKey
}

func (j *JWK) GetRSAKey() (*rsa.PrivateKey, error) {
	if j.privateRSAKey == nil {
		// parse the private key once
		block, _ := pem.Decode([]byte(j.PrivateKey))
		if block == nil || block.Type != rsaPrivateKeyType {
			log.Fatal("failed to decode PEM block containing private key")
		}
		pkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		j.privateRSAKey = pkey
	}
	return j.privateRSAKey, nil
}

// InsertJWK adds a JWK par
func (d *DAO) InsertJWKS(ctx context.Context, j *JWK) error {
	result, err := d.db.ExecContext(ctx, insertJWKSQL, j.PrivateKey, j.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to insert jwks: %+w", err)
	}
	countRows, _ := result.RowsAffected()
	if countRows != 1 {
		return fmt.Errorf("unexpected error inserting jwks: %+v", j)
	}
	return nil
}

func (d *DAO) GetJWKS(ctx context.Context) ([]*JWK, error) {
	rows, err := d.db.QueryContext(ctx, selectJWKSSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to select jwks: %w", err)
	}
	defer rows.Close()

	var results []*JWK
	for rows.Next() {
		j := new(JWK)
		if err := rows.Scan(&j.PrivateKey, &j.PublicKey); err != nil {
			return nil, fmt.Errorf("failed to scan jwk: %w", err)
		}
		results = append(results, j)
	}

	return results, nil
}
