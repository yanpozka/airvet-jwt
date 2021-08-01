package api

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

const (
	jwtIssuer = "Secure Company"
)

var (
	errInvalidJWT error = fmt.Errorf("invalid jwt")
)

type userClaims struct {
	jwt.Claims
	Email string `json:"email"`
}

func newJWT(email string, priKey *rsa.PrivateKey, expireAt time.Time) (string, error) {
	claims := userClaims{
		Email: email,
		Claims: jwt.Claims{
			Issuer: jwtIssuer,
			Expiry: jwt.NewNumericDate(expireAt),
		},
	}
	opts := jose.SignerOptions{}
	opts.WithType("JWT")
	// TODO: have a valid domain so jwt.io can verify our JWT :)
	opts.WithHeader("jku", "http://localhost:8080/jwks")

	signKey := jose.SigningKey{
		Algorithm: jose.RS256,
		Key:       priKey,
	}

	signer, err := jose.NewSigner(signKey, &opts)
	if err != nil {
		return "", err
	}

	return jwt.Signed(signer).
		Claims(claims).
		CompactSerialize()
}

func parseJWT(signedJWT string, pubKey *rsa.PublicKey) (*userClaims, error) {
	token, err := jwt.ParseSigned(signedJWT)
	if err != nil {
		return nil, errInvalidJWT
	}

	claims := new(userClaims)
	if err := token.Claims(pubKey, claims); err != nil {
		return nil, errInvalidJWT
	}

	err = claims.Validate(jwt.Expected{
		Issuer: jwtIssuer,
		Time:   time.Now(),
	})
	if err != nil {
		if err == jwt.ErrExpired {
			return nil, errInvalidJWT
		}

		return nil, errInvalidJWT
	}

	return claims, nil
}
