package controller

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/gonzabosio/res-manager/controller/handlers"
	middlewares "github.com/gonzabosio/res-manager/controller/middleware"
)

func Routing() *chi.Mux {
	r := chi.NewRouter()
	h, err := handlers.NewHandler()
	if err != nil {
		log.Fatal(err)
	}
	r.Use(middleware.Logger)
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONT_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	}))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Resources Manager"))
	})

	r.Group(func(r chi.Router) {
		r.Use(middlewares.OAuthMiddleware)
		r.Post("/team", h.CreateTeam)
		r.Post("/join-team", h.VerifyTeamByName)
		r.Get("/team", h.GetTeams)
		r.Patch("/team", h.ModifyTeam)
		r.Delete("/team/{team-id}", h.DeleteTeam)

		r.Post("/project", h.CreateProject)
		r.Get("/project/{team-id}", h.GetProjectsByTeamID)
		r.Patch("/project", h.ModifyProject)
		r.Delete("/project/{project-id}", h.DeleteProject)

		r.Post("/section", h.CreateSection)
		r.Get("/section/{project-id}", h.GetSectionsByProjectID)
		r.Put("/section", h.ModifySection)
		r.Delete("/section/{section-id}", h.DeleteSection)

		r.Post("/resource", h.CreateResource)
		r.Get("/resource/{section-id}", h.GetResourcesBySectionID)
		r.Patch("/resource", h.ModifyResource)
		r.Delete("/resource/{resource-id}", h.DeleteResource)

		r.Post("/login-user", h.LoginUser)
		r.Post("/user", h.RegisterUser) // if user already exists retrieve info
		// r.Get("/user", h.GetUsers)
		r.Patch("/user", h.ModifyUser)
		r.Delete("/user/{user-id}", h.DeleteUser)

		r.Post("/participant", h.AddParticipant)
		r.Get("/participant/{team-id}", h.GetParticipants)
		// update user role
		r.Delete("/participant/{user-id}/{team-id}", h.DeleteParticipant)

		r.Post("/csv", h.UploadCSV)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Get("/google_login", h.GoogleLoginHandler)
		r.Get("/google_callback", h.GoogleCallbackHandler)
	})
	return r
}
