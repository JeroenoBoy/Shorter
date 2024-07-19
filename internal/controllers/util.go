package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JeroenoBoy/Shorter/api"
	"github.com/JeroenoBoy/Shorter/internal/authentication"
	"github.com/JeroenoBoy/Shorter/internal/models"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

func WrapFunc(handler HandlerFunc) http.HandlerFunc {
	return WrapHandler(handler)
}

func WrapHandler(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler.ServeHTTP(w, r)
		if err == nil {
			return
		}

		err = api.WriteError(w, err)
		if err != nil {
			panic(err)
		}
	}
}

func ParseUserId(r *http.Request, paramName string, allowMe bool, allowAllPermissions models.Permissions) (models.UserId, error) {
	param, ok := r.Context().Value(paramName).(string)
	if !ok {
		return 0, fmt.Errorf("ParseUser: param '%v' does not exist", paramName)
	}

	user, hasUser := authentication.GetUser(r)

	if param == "@me" && allowMe {
		if !hasUser {
			return 0, api.ErrorNoPermissions
		}
		return user.Id, nil
	}

	if !hasUser && allowAllPermissions != 0 {
		return 0, api.ErrorNoPermissions
	} else if hasUser && !user.Permissions.HasAll(allowAllPermissions) {
		return 0, api.ErrorNoPermissions
	}

	id, err := strconv.Atoi(param)
	if err != nil {
		return 0, api.ErrorBadRequest
	}

	return models.UserId(id), nil
}
