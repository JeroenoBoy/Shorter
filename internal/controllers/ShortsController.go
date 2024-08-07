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
	r.Delete("/{id}", WrapAjaxFunc(c.deleteLink))
	r.Put("/{id}", WrapAjaxFunc(c.putLink))
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
		if strings.ToLower(ln) == "d" || strings.ToLower(ln) == "static" {
			return api.NewApiError(http.StatusBadRequest, "link 'd' and 'static' are reserved")
		}
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
	_, link, err := getUserAndLink(c.store, r)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return view.ShortRow(link).Render(r.Context(), w)
}

func (c *shortsController) getEditForm(w http.ResponseWriter, r *http.Request) error {
	_, link, err := getUserAndLink(c.store, r)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return view.ShortRowEdit(link).Render(r.Context(), w)
}

func (c *shortsController) putLink(w http.ResponseWriter, r *http.Request) error {
	_, link, err := getUserAndLink(c.store, r)
	if err != nil {
		return err
	}

	r.ParseForm()

	ln := r.Form.Get("link")
	if strings.ToLower(ln) == "d" || strings.ToLower(ln) == "static" {
		return api.NewApiError(http.StatusBadRequest, "link 'd' and 'static' are reserved")
	}
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

	target := r.Form.Get("target")
	if len(target) < 5 || (!strings.HasPrefix(target, "https://") && !strings.HasPrefix(target, "http://")) {
		return api.NewApiError(http.StatusBadRequest, "target must start with http:// or https://")
	}

	if utf8.RuneCountInString(target) > 512 {
		return api.NewApiError(http.StatusBadRequest, "target must not be longer than 512 characters")
	}

	updateReq := api.UpdateLinkRequest{
		Link:   ln,
		Target: target,
	}

	link, err = c.store.UpdateLink(link.Id, updateReq)
	if err != nil {
		if errors.Is(err, datastore.ErrorDuplicateKey) {
			return api.ErrorDuplicateKey
		}
		return err
	}

	w.WriteHeader(http.StatusOK)
	return view.ShortRow(link).Render(r.Context(), w)
}

func (c *shortsController) deleteLink(w http.ResponseWriter, r *http.Request) error {
	_, link, err := getUserAndLink(c.store, r)
	if err != nil {
		return err
	}

	err = c.store.DeleteLink(link.Id)
	if err != nil {
		if errors.Is(err, datastore.ErrorLinkNotFound) {
			return api.ErrorResourceNotFound
		}
		return err
	}

	w.WriteHeader(200)
	w.Write([]byte{})
	return nil
}

func getUserAndLink(store datastore.Datastore, r *http.Request) (user models.User, link models.ShortLink, err error) {
	user, ok := authentication.GetUser(r)
	if !ok {
		return models.User{}, models.ShortLink{}, api.ErrorNotAuthenticated
	}

	linkId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return models.User{}, models.ShortLink{}, api.ErrorResourceNotFound
	}

	link, err = store.GetLink(models.LinkId(linkId))
	if err != nil {
		if errors.Is(err, datastore.ErrorLinkNotFound) {
			return models.User{}, models.ShortLink{}, api.ErrorResourceNotFound
		}
		return models.User{}, models.ShortLink{}, err
	}

	if link.Owner != user.Id && !user.Permissions.HasAll(models.PermissionsManageShorts) {
		return models.User{}, models.ShortLink{}, api.ErrorNoPermissions
	}

	return user, link, nil
}
