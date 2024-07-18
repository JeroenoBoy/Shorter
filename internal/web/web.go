package web

import (
	"net/http"

	"github.com/JeroenoBoy/Shorter/internal/authentication"
	"github.com/JeroenoBoy/Shorter/internal/controllers"
	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RunServer(datastore datastore.Datastore) error {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(middleware.Heartbeat("/ping"))
		r.Use(authentication.MiddilewareProvideUser)

		r.Mount("/shorts", controllers.NewShortsController(datastore).Router())
	})

	return http.ListenAndServe(":3000", r)
}
