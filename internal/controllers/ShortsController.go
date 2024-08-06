package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/JeroenoBoy/Shorter/api"
	"github.com/JeroenoBoy/Shorter/internal/authentication"
	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/JeroenoBoy/Shorter/internal/models"
	"github.com/JeroenoBoy/Shorter/view"
	"github.com/go-chi/chi/v5"
)

type shortsController struct {
	store datastore.Datastore
}

func NewShortController(store datastore.Datastore) *shortsController {
	return &shortsController{
		store: store,
	}
}

func (c *shortsController) Router() chi.Router {
	r := chi.NewRouter()
	r.Post("/", WrapAjaxFunc(c.newLink))
	r.Get("/{id}", WrapAjaxFunc(c.getLink))
	r.Get("/{id}/edit", WrapAjaxFunc(c.getEditForm))
	return r
}

func (c *shortsController) newLink(w http.ResponseWriter, r *http.Request) error {
	user, ok := authentication.GetUser(r)
	if !ok {
		return api.ErrorNotAuthenticated
	}

	err := r.ParseForm()
	if err != nil {
		return api.ErrorBadRequest
	}

	if !r.Form.Has("target") {
		return api.ErrorBadRequest
	}

	target := strings.TrimSpace(r.Form.Get("target"))
	if len(target) < 5 || (!strings.HasPrefix(target, "https://") && !strings.HasPrefix(target, "http://")) {
		return api.NewApiError(http.StatusBadRequest, "target must start with http:// or https://")
	}

	if utf8.RuneCountInString(target) > 512 {
		return api.NewApiError(http.StatusBadRequest, "target must not be longer than 512 characters")
	}

	var link *string
	if r.Form.Has("link") {
		ln := r.Form.Get("link")
		if utf8.RuneCountInString(ln) > 24 {
			return api.NewApiError(http.StatusBadRequest, "link must not be longer than 24 characters")
		}
		if m, err := regexp.Match("[a-zA-Z0-9_\\-]*", []byte(ln)); err != nil || !m {
			if err != nil {
				return err
			} else {
				return api.NewApiError(http.StatusBadRequest, "link must follow 'a-zA-Z0-9_-'")
			}
		}
		link = &ln
	} else {
		link = nil
	}

	createdLn, err := c.store.CreateLink(user.Id, link, target)
	if err != nil {
		if errors.Is(err, datastore.ErrorDuplicateKey) {
			return api.ErrorDuplicateKey
		}
		return err
	}

	w.WriteHeader(http.StatusOK)
	return view.ShortRow(createdLn).Render(r.Context(), w)
}

func (c *shortsController) getLink(w http.ResponseWriter, r *http.Request) error {
	user, ok := authentication.GetUser(r)
	if !ok {
		return api.ErrorNotAuthenticated
	}

	linkId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return api.ErrorResourceNotFound
	}

	link, err := c.store.GetLink(models.LinkId(linkId))
	if err != nil {
		if errors.Is(err, datastore.ErrorLinkNotFound) {
			return api.ErrorResourceNotFound
		}
		return err
	}
	if link.Owner != user.Id && !user.Permissions.HasAll(models.PermissionsManageShorts) {
		return api.ErrorNoPermissions
	}
	w.WriteHeader(http.StatusOK)
	return view.ShortRow(link).Render(r.Context(), w)
}

func (c *shortsController) getEditForm(w http.ResponseWriter, r *http.Request) error {
	user, ok := authentication.GetUser(r)
	if !ok {
		return api.ErrorNotAuthenticated
	}

	linkId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return api.ErrorResourceNotFound
	}

	link, err := c.store.GetLink(models.LinkId(linkId))
	if err != nil {
		if errors.Is(err, datastore.ErrorLinkNotFound) {
			return api.ErrorResourceNotFound
		}
		return err
	}

	if link.Owner != user.Id && !user.Permissions.HasAll(models.PermissionsManageShorts) {
		return api.ErrorNoPermissions
	}

	w.WriteHeader(http.StatusOK)
	return view.ShortRowEdit(link).Render(r.Context(), w)
}
