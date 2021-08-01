package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/square/go-jose/v3"
)

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
		rsaKey, _ := jwkDB.GetRSAKey()

		jwk := jose.JSONWebKey{
			Key:       &rsaKey.PublicKey,
			Algorithm: "RS256",
			Use:       "sig",
		}

		jwks.Keys = append(jwks.Keys, jwk)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jwks)
}
