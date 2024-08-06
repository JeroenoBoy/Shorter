package controllers

import (
	"errors"
	"net/http"

	"github.com/JeroenoBoy/Shorter/api"
	"github.com/JeroenoBoy/Shorter/internal/authentication"
	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/JeroenoBoy/Shorter/view"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type loginController struct {
	auth  *authentication.JWTAuthentication
	store datastore.Datastore
}

func NewLoginController(store datastore.Datastore, auth *authentication.JWTAuthentication) *loginController {
	return &loginController{
		auth:  auth,
		store: store,
	}
}

func (c *loginController) Router() chi.Router {
	r := chi.NewRouter()
	r.Use(c.auth.MiddilewareProvideUser)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := authentication.GetUser(r); ok {
				http.Redirect(w, r, "/d", http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
	r.Get("/", middleware.SetHeader("HX-Redirect", "/d/login")(templ.Handler(view.LoginPage())).ServeHTTP)
	r.Post("/", WrapAjaxFunc(c.LoginUser))
	return r
}

func (c *loginController) LoginUser(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return api.ErrorBadRequest
	}

	if !r.Form.Has("username") || !r.Form.Has("password") {
		return api.ErrorBadRequest
	}

	usr := r.Form.Get("username")
	pwd := r.Form.Get("password")

	user, err := c.store.FindUserByName(usr)
	if err != nil {
		if errors.Is(err, datastore.ErrorUserNotFound) {
			w.WriteHeader(http.StatusOK)
			return view.LoginFailedMessage("No matching user found").Render(r.Context(), w)
		}
		return err
	}

	err = user.Passwd.Compare(pwd)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return view.LoginFailedMessage("No matching user found").Render(r.Context(), w)
	}

	cookie, err := c.auth.CreateCookie(user)
	if err != nil {
		return err
	}

	http.SetCookie(w, &cookie)
	view.HtmxRedirect(w, "/d")
	return nil
}
