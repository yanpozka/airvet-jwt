package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/yanpozka/airvet-jwt/dao"
)

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
		if err == dao.UserNotFoundErr {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		} else {
			log.Printf("Error getting user from db: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	log.Println(*user)

	http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
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
