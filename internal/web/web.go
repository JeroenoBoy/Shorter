package web

import (
	"net/http"

	"github.com/JeroenoBoy/Shorter/api"
	"github.com/JeroenoBoy/Shorter/internal/authentication"
	"github.com/JeroenoBoy/Shorter/internal/controllers"
	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/JeroenoBoy/Shorter/view"
	"github.com/a-h/templ"
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

	fs := http.FileServer(http.Dir("static"))

	r.Handle("/static/*", http.StripPrefix("/static", fs))
	r.Mount("/login", controllers.NewLoginController(datastore, jwtAuth).Router())

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(jwtAuth.MiddilewareProvideUser)

		r.Mount("/shorts", controllers.NewShortsController(datastore).Router())
	})

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(jwtAuth.MiddilewareProvideUser)
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, ok := authentication.GetUser(r); !ok {
					println("no user tho")
					http.Redirect(w, r, "/login", http.StatusFound)
				} else {
					println("with user tho")
					next.ServeHTTP(w, r)
				}
			})
		})

		r.Get("/", templ.Handler(view.IndexPage()).ServeHTTP)

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			view.WriteErrorPage(w, r.Context(), api.NewApiError(http.StatusNotFound, "Page could not be found :("))
		})
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			view.WriteError(w, r.Context(), api.NewApiError(http.StatusMethodNotAllowed, "Api Error: Method Not Allowed"))
		})
	})

	return r
}
