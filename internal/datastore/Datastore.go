package datastore

import "github.com/JeroenoBoy/Shorter/internal/models"

type Datastore interface {
	GetUser(id models.UserId) (models.User, error)
	GetUsers() ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(user models.UserId) error

	CreateShort(models.ShortLink) error
	DeleteShort(models.ShortId) error
	GetShort(id models.ShortId) (models.ShortLink, error)
	GetUserShorts(userId models.UserId) ([]models.ShortLink, error)
	GetAllShorts() ([]models.ShortLink, error)
}
