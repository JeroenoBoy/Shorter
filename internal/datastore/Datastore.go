package datastore

import (
	"errors"

	"github.com/JeroenoBoy/Shorter/api"
	"github.com/JeroenoBoy/Shorter/internal/models"
)

type Datastore interface {
	GetUsers(request models.PageRequest) (models.PaginatedUsers, error)
	GetUser(id models.UserId) (models.User, error)
	FindUserByName(name string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(user models.UserId) error

	GetLink(id models.LinkId) (models.ShortLink, error)
	GetAllLinks(request models.PageRequest) (models.PaginatedLinks, error)
	GetUserLinks(owner models.UserId) ([]models.ShortLink, error)
	CreateLink(owner models.UserId, link *string, target string) (models.ShortLink, error)
	UpdateLink(id models.LinkId, updateReq api.UpdateLinkRequest) (models.ShortLink, error)
	DeleteLink(id models.LinkId) error

	GetLinkTargetAndIncreaseRedirects(link string) (string, error)
}

var (
	ErrorUserNotFound     = errors.New("user not found")
	ErrorLinkNotFound     = errors.New("link not found")
	ErrorNoDataReceived   = errors.New("no data received")
	ErrorDuplicateKey     = errors.New("duplicate key")
	ErrorInPreperation    = errors.New("error while preparing statement")
	ErrorInRequest        = errors.New("error while trying to execute statement")
	ErrorInScan           = errors.New("error while scanning")
	ErrorInvalidPage      = errors.New("page number must be zero or greater")
	ErrorInvalidItemCount = errors.New("must have atleast 1 item per page")
)
