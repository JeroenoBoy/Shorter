package web

import (
	"errors"
	"net/http"

	"github.com/JeroenoBoy/Shorter/api"
	"github.com/JeroenoBoy/Shorter/internal/authentication"
	"github.com/JeroenoBoy/Shorter/internal/controllers"
	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/JeroenoBoy/Shorter/internal/models"
	"github.com/JeroenoBoy/Shorter/view"
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
	r.Mount("/d/login", controllers.NewLoginController(datastore, jwtAuth).Router())

	r.Route("/d", func(r chi.Router) {
		r.Use(middleware.NoCache)
		r.Use(jwtAuth.MiddilewareProvideUser)
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if _, ok := authentication.GetUser(r); !ok {
					api.Redirect(w, r, "/d/login")
				} else {
					next.ServeHTTP(w, r)
				}
			})
		})

		r.Get("/", controllers.WrapPageFunc(IndexPage(datastore)))
		r.Mount("/shorts", controllers.NewUserShortController(datastore).Router())

		r.Route("/admin", func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					user, ok := authentication.GetUser(r)
					if !ok {
						api.Redirect(w, r, "/d")
					} else if !user.Permissions.HasAny(models.PermissionsAnyDashboardAccess) {
						api.Redirect(w, r, "/d")
					} else {
						next.ServeHTTP(w, r)
					}
				})
			})

			r.Get("/", controllers.WrapPageFunc(AdminPage(datastore)))
			r.Get("/links", controllers.WrapPageFunc(AdminLinksPage(datastore)))
			r.Mount("/users", controllers.NewUserController(datastore).Router())
			r.Get("/settings", controllers.WrapPageFunc(AdminSettingsPage(datastore)))
		})

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			err := api.NewApiError(http.StatusNotFound, "Page could not be found :(")
			if api.IsHTMXRequest(r) {
				view.ErrorNotification(w, r.Context(), err)
			} else {
				view.WriteErrorPage(w, r.Context(), err)
			}
		})
		r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
			err := api.NewApiError(http.StatusMethodNotAllowed, "Api Error: Method Not Allowed")
			if api.IsHTMXRequest(r) {
				view.ErrorNotification(w, r.Context(), err)
			} else {
				view.WriteErrorPage(w, r.Context(), err)
			}
		})
	})

	r.Get("/{link}", controllers.WrapPageFunc(RedirectRoute(datastore)))

	return r
}

func IndexPage(store datastore.Datastore) controllers.HandlerFunc {
	return controllers.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		user, ok := authentication.GetUser(r)
		if !ok {
			return api.NewApiError(http.StatusUnauthorized, "You must be logged in to view this")
		}

		links, err := store.GetUserLinks(user.Id)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)
		return view.ShortsPage(user, links).Render(r.Context(), w)
	})
}

func RedirectRoute(store datastore.Datastore) controllers.HandlerFunc {
	return controllers.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		link := r.PathValue("link")
		if len(link) == 0 || len(link) > 24 {
			return api.ErrorResourceNotFound
		}

		target, err := store.GetLinkTargetAndIncreaseRedirects(link)
		if err != nil {
			if errors.Is(err, datastore.ErrorLinkNotFound) {
				return api.ErrorResourceNotFound
			}
			return err
		}

		http.Redirect(w, r, target, http.StatusFound)
		return nil
	})
}
