package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/square/go-jose/v3"
)

// we could also use RS512
const jwkAlgo = "RS256"

func (a *API) getJWKS(w http.ResponseWriter, req *http.Request) {
	jwksDB, err := a.db.GetJWKS(req.Context())
	if err != nil {
		log.Println("Error gettings JWKS: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jwks := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{},
	}

	for _, jwkDB := range jwksDB {
		publicKey, err := jwkDB.GetRSAPublicKey()
		if err != nil {
			log.Printf("Error gettings JWK public key: %v", err)
			continue
		}

		jwk := jose.JSONWebKey{
			Key:       publicKey,
			Algorithm: jwkAlgo,
			Use:       "sig",
		}

		jwks.Keys = append(jwks.Keys, jwk)
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jwks)
}
