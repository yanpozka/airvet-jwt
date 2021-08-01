package dao

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
)

const (
	keySize = 2048

	rsaPrivateKeyType = "RSA PRIVATE KEY"
	rsaPublicKeyType  = "PUBLIC KEY"
)

// GeneratePrivatePublicKeyPair returns a private and public key pair
// TODO: move this function to its own pkg
func GeneratePrivatePublicKeyPair() (string, string, error) {
	privatekey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return "", "", err
	}
	publickey := &privatekey.PublicKey

	privatePem := new(strings.Builder)
	{
		privateKeyBytes := x509.MarshalPKCS1PrivateKey(privatekey)
		privateKeyBlock := &pem.Block{
			Type:  rsaPrivateKeyType,
			Bytes: privateKeyBytes,
		}

		if err = pem.Encode(privatePem, privateKeyBlock); err != nil {
			return "", "", fmt.Errorf("error when encode private pem: %w", err)
		}
	}

	publicPem := new(strings.Builder)
	{
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
		if err != nil {
			return "", "", err
		}
		publicKeyBlock := &pem.Block{
			Type:  rsaPublicKeyType,
			Bytes: publicKeyBytes,
		}

		if err = pem.Encode(publicPem, publicKeyBlock); err != nil {
			return "", "", fmt.Errorf("error when encode private pem: %w", err)
		}
	}

	return privatePem.String(), publicPem.String(), nil
}
