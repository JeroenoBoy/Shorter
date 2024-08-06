package view

import (
	"context"
	"net/http"

	"github.com/JeroenoBoy/Shorter/api"
)

func WriteErrorPage(w http.ResponseWriter, ctx context.Context, err error) error {
	if apiErr, ok := err.(api.ApiError); ok {
		w.WriteHeader(apiErr.StatusCode)
	} else {
		defer panic(err)
		w.WriteHeader(500)
	}
	return ErrorPage(err).Render(ctx, w)
}

func HtmxRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("HX-Redirect", url)
	w.WriteHeader(http.StatusOK)
	w.Write(make([]byte, 0))
}
