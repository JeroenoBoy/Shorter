package api

import "net/http"

type UpdateLinkRequest struct {
	Link   string `json:"link"`
	Target string `json:"target"`
}

func IsHTMXRequest(r *http.Request) bool {
	htmxHeader := r.Header.Get("HX-Request")
	return htmxHeader == "true"
}

func Redirect(w http.ResponseWriter, r *http.Request, url string) {
	if IsHTMXRequest(r) {
		w.Header().Set("HX-Redirect", url)
		w.WriteHeader(http.StatusOK)
		w.Write(make([]byte, 0))
	} else {
		http.Redirect(w, r, url, http.StatusFound)
	}
}
