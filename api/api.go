package api

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

func WriteResponse(w http.ResponseWriter, code int, data any) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(json)
	return err
}

func WriteMessage(w http.ResponseWriter, code int, message string) error {
	return WriteResponse(w, code, Message{StatusCode: code, Message: message})
}

func WriteOk(w http.ResponseWriter) error {
	return WriteMessage(w, http.StatusOK, "Ok")
}
