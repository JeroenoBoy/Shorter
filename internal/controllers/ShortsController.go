package controllers

import (
	"net/http"

	"github.com/JeroenoBoy/Shorter/api"
	"github.com/JeroenoBoy/Shorter/internal/authentication"
	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/JeroenoBoy/Shorter/internal/models"
	"github.com/go-chi/chi/v5"
)

type shortsController struct {
	store datastore.Datastore
}

func NewShortsController(datastore datastore.Datastore) *shortsController {
	return &shortsController{
		store: datastore,
	}
}

func (c *shortsController) Router() chi.Router {
	r := chi.NewRouter()
	r.Use(authentication.MiddlewareIsAuthenticated)
	r.Get("/", authentication.MiddlewareHasPermissions(models.PermissionsManageShorts)(WrapApiFunc(c.getAllShorts)).ServeHTTP)
	r.Get("/{userid}", WrapApiFunc(c.getShortsForUser))
	return r
}

func (c *shortsController) getShortsForUser(w http.ResponseWriter, r *http.Request) error {
	uid, err := ParseUserId(r, "userId", true, models.PermissionsManageShorts)
	if err != nil {
		return err
	}

	data, err := c.store.GetUserShorts(uid)
	if err != nil {
		return err
	}

	return api.WriteResponse(w, http.StatusOK, data)
}

func (c *shortsController) getAllShorts(w http.ResponseWriter, r *http.Request) error {
	data, err := c.store.GetAllShorts()
	if err != nil {
		return err
	}

	return api.WriteResponse(w, http.StatusOK, data)
}
