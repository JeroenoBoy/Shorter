package authentication

import (
	"net/http"

	"github.com/JeroenoBoy/Shorter/api"
	"github.com/JeroenoBoy/Shorter/internal/models"
)

type authCtxKey struct {
	name string
}

var (
	CtxKeyClaims = &authCtxKey{"Claims"}
	CtxKeyUser   = &authCtxKey{"User"}
)

func MiddlewareIsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetUser(r)
		if !ok {
			api.WriteError(w, api.ErrorNotAuthenticated)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func MiddlewareIsNotAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := GetUser(r)
		if ok {
			api.WriteError(w, api.ErrorNoPermissions)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func MiddlewareHasPermissions(permissions models.Permissions) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUser(r)
			if !ok {
				api.WriteError(w, api.ErrorNotAuthenticated)
				return
			}

			if user.Permissions.HasAll(permissions) {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUser(r *http.Request) (user models.User, ok bool) {
	user, ok = r.Context().Value(CtxKeyClaims).(models.User)
	return user, ok
}
