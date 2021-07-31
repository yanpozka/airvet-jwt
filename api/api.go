package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/yanpozka/airvet-jwt/dao"
)

const jwtAddExpiry = 30 * 24 * time.Hour // 1 month

type API struct {
	db *dao.DAO
}

type userIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAPI(db *dao.DAO) *API {
	return &API{
		db: db,
	}
}

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

	rsaKey, err := user.GetRSAKey()
	if err != nil {
		log.Printf("Error decoding private key: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	jwt, err := newJWT(user.Email, rsaKey, time.Now().Add(jwtAddExpiry))
	if err != nil {
		log.Printf("Error generating jwt: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := map[string]string{
		"jwt": jwt,
	}
	json.NewEncoder(w).Encode(resp)
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
