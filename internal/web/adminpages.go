package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JeroenoBoy/Shorter/api"
	"github.com/JeroenoBoy/Shorter/internal/authentication"
	"github.com/JeroenoBoy/Shorter/internal/controllers"
	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/JeroenoBoy/Shorter/internal/models"
	"github.com/JeroenoBoy/Shorter/view/adminview"
)

func AdminPage(store datastore.Datastore) controllers.HandlerFunc {
	return controllers.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		user, ok := authentication.GetUser(r)
		if !ok {
			return api.NewApiError(http.StatusUnauthorized, "You must be logged in to view this")
		}

		if user.Permissions.HasAll(models.PermissionsManageShorts) {
			api.Redirect(w, r, "/d/admin/links")
			return nil
		} else if user.Permissions.HasAll(models.PermissionsManageUsers) {
			api.Redirect(w, r, "/d/admin/users")
			return nil
		} else if user.Permissions.HasAll(models.PermissionsManageServer) {
			api.Redirect(w, r, "/d/admin/settings")
			return nil
		}

		return fmt.Errorf("user '%v'(%v) has no access the admin panel", user.Name, user.Id)
	})
}

func AdminLinksPage(store datastore.Datastore) controllers.HandlerFunc {
	return controllers.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		user, ok := authentication.GetUser(r)
		if !ok {
			return api.NewApiError(http.StatusUnauthorized, "You must be logged in to view this")
		}

		pageParam := r.URL.Query().Get("page")
		if len(pageParam) == 0 {
			pageParam = "1"
		}

		pageNum, err := strconv.Atoi(pageParam)
		if err != nil {
			return api.ErrorBadRequest
		}

		pageRequest := models.PageRequest{
			Page:         pageNum,
			ItemsPerPage: 15,
		}

		links, err := store.GetAllLinks(pageRequest)
		if err != nil {
			return err
		}

		if api.IsHTMXRequest(r) {
			w.WriteHeader(http.StatusOK)
			return adminview.LinksTable(links).Render(r.Context(), w)
		} else {
			w.WriteHeader(http.StatusOK)
			return adminview.LinksPage(user, links).Render(r.Context(), w)
		}

	})
}

func AdminSettingsPage(store datastore.Datastore) controllers.HandlerFunc {
	return controllers.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		user, ok := authentication.GetUser(r)
		if !ok {
			return api.NewApiError(http.StatusUnauthorized, "You must be logged in to view this")
		}

		w.WriteHeader(http.StatusOK)
		return adminview.SettingsPage(user).Render(r.Context(), w)
	})
}
