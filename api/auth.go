package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/yanpozka/airvet-jwt/dao"
)

const (
	jwtExpiration = 30 * 24 * time.Hour // 1 month

	authorizationHeader = "Authorization"
)

func (a *API) auth(w http.ResponseWriter, req *http.Request) {
	u, err := readUserIn(req)
	if err != nil {
		log.Println("Error reading user input: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	user, err := a.db.GetUserByEmailPasswd(req.Context(), u.Email, u.Password)
	if err != nil {
		if err == dao.ErrUserNotFound {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		} else {
			log.Printf("Error getting user from db: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	jwks, err := a.db.GetJWKS(req.Context())
	if len(jwks) == 0 || err != nil {
		log.Printf("We don't have a JWKS, error?: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jwk := jwks[0] // always get the first one
	rsaKey, err := jwk.GetRSAPrivateKey()
	if err != nil {
		log.Printf("Error decoding private key: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jwt, err := newJWT(user.Email, rsaKey, time.Now().Add(jwtExpiration))
	if err != nil {
		log.Printf("Error generating jwt: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp := map[string]string{
		"jwt": jwt,
	}
	json.NewEncoder(w).Encode(resp)
}
