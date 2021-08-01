package api

import "github.com/yanpozka/airvet-jwt/dao"

// API represents the whole api
type API struct {
	db *dao.DAO
}

// NewAPI creates a new API
func NewAPI(db *dao.DAO) *API {
	return &API{
		db: db,
	}
}
