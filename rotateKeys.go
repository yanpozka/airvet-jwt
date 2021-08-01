package main

import (
	"context"
	"log"
	"time"

	"github.com/yanpozka/airvet-jwt/dao"
)

const dbPath = "users.db"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	d, err := dao.NewDAO(dbPath)
	if err != nil {
		log.Panic(err)
	}

	privatekey, publickey, err := dao.GeneratePrivatePublicKeyPair()
	if err != nil {
		log.Panic(err)
	}
	expTime := time.Now().Add(dao.JWKExpiration)
	jwk := &dao.JWK{
		PrivateKey: privatekey,
		PublicKey:  publickey,
		ExpiresAt:  expTime.Unix(),
	}
	if err := d.InsertJWK(context.Background(), jwk); err != nil {
		log.Panic(err)
	}

	log.Printf("Added a new JWK, will expire at: %v", expTime)
}
