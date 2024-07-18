package datastore

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/JeroenoBoy/Shorter/internal/models"
	_ "github.com/lib/pq"
)

type PostgresOptions struct {
	Database string
	UserName string
	Password string
	Ip       string
	Port     int
	UseSSL   bool
}

type postgresStore struct {
	*sql.DB
}

func CreatePostgrsStore(options PostgresOptions) (Datastore, error) {
	if len(options.Ip) == 0 {
		options.Ip = "localhost"
	}
	if options.Port == 0 {
		options.Port = 5432
	}

	var sslmode string
	if options.UseSSL {
		sslmode = "enabled"
	} else {
		sslmode = "disabled"
	}

	connStr := fmt.Sprint("postgres://%v:%v@%v:%v/%v?sslmode=%v", options.UserName, options.Password, options.Ip, strconv.Itoa(options.Port), options.Database, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &postgresStore{
		DB: db,
	}, nil
}

func (s *postgresStore) GetUser(id models.UserId) (models.User, error) {
	panic("not implemented")
}

func (s *postgresStore) GetUsers() ([]models.User, error) {
	panic("not implemented")
}

func (s *postgresStore) CreateUser(user models.User) (models.User, error) {
	panic("not implemented")
}

func (s *postgresStore) UpdateUser(user models.User) (models.User, error) {
	panic("not implemented")
}

func (s *postgresStore) DeleteUser(user models.UserId) error {
	panic("not implemented")
}

func (s *postgresStore) CreateShort(models.ShortLink) error {
	panic("not implemented")
}

func (s *postgresStore) DeleteShort(models.ShortId) error {
	panic("not implemented")
}

func (s *postgresStore) GetShort(id models.ShortId) (models.ShortLink, error) {
	panic("not implemented")
}

func (s *postgresStore) GetUserShorts(userId models.UserId) ([]models.ShortLink, error) {
	panic("not implemented")
}

func (s *postgresStore) GetAllShorts() ([]models.ShortLink, error) {
	panic("not implemented")
}
