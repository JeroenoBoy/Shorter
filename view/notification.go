package view

import (
	"context"
	"net/http"
	"strconv"

	"github.com/JeroenoBoy/Shorter/api"
)

type NotificationType string
type Notification struct {
	Header string
	Body   string
	Type   NotificationType
}

const (
	NotificationTypeError   = NotificationType("error")
	NotificationTypeWarning = NotificationType("warning")
	NotificationTypeNotify  = NotificationType("notify")
)

var NotificationInternalServerError = CreateNotification(NotificationTypeError, "Internal Server Error!", "Whoopsies... Something went wrong!")

func CreateNotification(nType NotificationType, header string, body string) Notification {
	return Notification{
		Type:   nType,
		Header: header,
		Body:   body,
	}
}

func SendNotification(w http.ResponseWriter, ctx context.Context, notification Notification) error {
	w.Header().Set("HX-Retarget", "#notifications-container")
	w.Header().Set("HX-Reswap", "beforeend")
	w.WriteHeader(http.StatusOK)
	return NotificationTemplate(notification).Render(ctx, w)
}

func ErrorNotification(w http.ResponseWriter, ctx context.Context, err error) error {
	if apiErr, ok := err.(api.ApiError); ok {
		return SendNotification(w, ctx, CreateNotification(NotificationTypeError, "Error "+strconv.Itoa(apiErr.StatusCode), apiErr.Message))
	} else {
		SendNotification(w, ctx, NotificationInternalServerError)
		panic(err)
	}
}
