package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gonzabosio/res-manager/controller/handlers"
)

func Routing() *chi.Mux {
	r := chi.NewRouter()
	h, err := handlers.NewHandler()
	if err != nil {
		log.Fatal(err)
	}
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{fmt.Sprintf("http://localhost%v", os.Getenv("FRONT_PORT"))},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
	}))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Resources Manager"))
	})
	r.Post("/team", h.CreateTeam)
	return r
}