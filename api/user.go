package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type userIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) getUser(w http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get(authorizationHeader)
	parts := strings.Split(authHeader, " ")
	if len(parts) < 2 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	signedJWT := parts[1]

	jwks, err := a.db.GetJWKS(req.Context())
	if len(jwks) == 0 || err != nil {
		log.Printf("We don't have a JWKS, error?: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jwk := jwks[0] // always get the first one
	rsaKey, err := jwk.GetRSAKey()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	uc, err := parseJWT(signedJWT, &rsaKey.PublicKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := a.db.GetUserByEmail(req.Context(), uc.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func readUserIn(req *http.Request) (*userIn, error) {
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	u := new(userIn)
	if err := json.Unmarshal(data, u); err != nil {
		return nil, err
	}

	return u, nil
}
