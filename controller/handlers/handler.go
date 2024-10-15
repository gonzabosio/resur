package handlers

import (
	database "github.com/gonzabosio/res-manager/model/db"
	"github.com/gonzabosio/res-manager/model/db/repository"
)

type Handler struct {
	RepositoryService *repository.DBService
}

func NewHandler() (*Handler, error) {
	h := new(Handler)
	db, err := database.NewDB()
	if err != nil {
		return nil, err
	}
	s := &repository.DBService{DB: db}
	h.RepositoryService = s
	return h, nil
}