package datastore

import (
	"errors"
	"github.com/JeroenoBoy/Shorter/internal/models"
)

type Datastore interface {
	GetUsers() ([]models.User, error)
	GetUser(id models.UserId) (models.User, error)
	FindUserByName(name string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(user models.UserId) error

	CreateLink(owner models.UserId, link *string, target string) (models.ShortLink, error)
	DeleteLink(id models.LinkId) error
	GetLink(id models.LinkId) (models.ShortLink, error)
	GetUserLinks(owner models.UserId) ([]models.ShortLink, error)
	GetAllLinks() ([]models.ShortLink, error)

	GetLinkTargetAndIncreaseRedirects(link string) (string, error)
}

var (
	ErrorUserNotFound  = errors.New("user not found")
	ErrorLinkNotFound  = errors.New("link not found")
	ErrorDuplicateKey  = errors.New("duplicate key")
	ErrorInPreperation = errors.New("error while preparing statement")
	ErrorInRequest     = errors.New("error while trying to execute statement")
)
