package web

import (
	"net/http"

	"github.com/JeroenoBoy/Shorter/internal/authentication"
	"github.com/JeroenoBoy/Shorter/internal/controllers"
	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewWebserver(jwtAuth *authentication.JWTAuthentication, datastore datastore.Datastore) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(jwtAuth.MiddilewareProvideUser)

		r.Mount("/shorts", controllers.NewShortsController(datastore).Router())
	})

	return r
}
