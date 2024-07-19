package datastore

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/JeroenoBoy/Shorter/internal/models"
	_ "github.com/lib/pq"
)

type PostgresOptions struct {
	Database string
	UserName string
	Password string
	Host     string
	Port     int
	SSLMode  string
}

type postgresStore struct {
	*sql.DB
}

func NewPostgresStore(options PostgresOptions) (Datastore, error) {
	if len(options.Host) == 0 {
		return nil, errors.New("options.Host was empty in PostgresOptions")
	}

	if len(options.Database) == 0 {
		return nil, errors.New("options.Database was empty in PostgresOptions")
	}

	connOptions := make([]string, 0, 8)
	connOptions = append(connOptions, "host="+options.Host)
	connOptions = append(connOptions, "dbname="+options.Database)

	if options.Port > 0 {
		connOptions = append(connOptions, "port="+strconv.Itoa(options.Port))
	}

	if len(options.UserName) > 0 {
		connOptions = append(connOptions, "username="+options.UserName)
	}

	if len(options.Password) > 0 {
		connOptions = append(connOptions, "password="+options.Password)
	}

	if len(options.SSLMode) > 0 {
		connOptions = append(connOptions, "sslmode="+options.SSLMode)
	}

	connStr := strings.Join(connOptions, " ")
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
