package controllers

import (
	"net/http"
	"strconv"

	"github.com/JeroenoBoy/Shorter/api"
	"github.com/JeroenoBoy/Shorter/internal/authentication"
	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/JeroenoBoy/Shorter/internal/models"
	"github.com/JeroenoBoy/Shorter/view/adminview"
	"github.com/go-chi/chi/v5"
)

type userController struct {
	store datastore.Datastore
}

func NewUserController(store datastore.Datastore) *userController {
	return &userController{
		store: store,
	}
}

func (c *userController) Router() chi.Router {
	r := chi.NewRouter()
	r.Use(authentication.MiddlewareHasPermissionsOrRedirect(models.PermissionsManageUsers, "/d/admin"))

	r.Get("/", WrapPageFunc(c.userPage))
	return r
}

func (c *userController) userPage(w http.ResponseWriter, r *http.Request) error {
	user, ok := authentication.GetUser(r)
	if !ok {
		return api.NewApiError(http.StatusUnauthorized, "You must be logged in to view this")
	}

	pageRequest := models.PageRequest{
		Page:         0,
		ItemsPerPage: 10,
	}

	query := r.URL.Query()
	if query.Has("page") {
		n, err := strconv.Atoi(query.Get("page"))
		if err != nil {
			return api.ErrorBadRequest
		}
		if n < 1 {
			return api.ErrorBadRequest
		}
		pageRequest.ItemsPerPage = n
	}

	paginatedUseres, err := c.store.GetUsers(pageRequest)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return adminview.UsersPage(user, paginatedUseres).Render(r.Context(), w)
}
