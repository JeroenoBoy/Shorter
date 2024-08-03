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

	CreateShort(models.ShortLink) error
	DeleteShort(models.ShortId) error
	GetShort(id models.ShortId) (models.ShortLink, error)
	GetUserShorts(userId models.UserId) ([]models.ShortLink, error)
	GetAllShorts() ([]models.ShortLink, error)
}

var (
	ErrorUserNotFound  = errors.New("user not found")
	ErrorInPreperation = errors.New("error while preparing statement")
	ErrorInRequest     = errors.New("error while trying to query in getUser")
)
