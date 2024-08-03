package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/JeroenoBoy/Shorter/internal/models"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
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

	store := &postgresStore{
		DB: db,
	}

	err = store.initdb()
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *postgresStore) GetUsers() ([]models.User, error) {
	panic("not implemented")
}

func (s *postgresStore) GetUser(id models.UserId) (models.User, error) {
	var empty models.User
	rows, err := s.Query("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return empty, errors.Join(ErrorInRequest, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return empty, ErrorUserNotFound
	}
	rows.Scan(&empty)

	return userFromQuery(rows)
}

func (s *postgresStore) FindUserByName(name string) (models.User, error) {
	var empty models.User

	rows, err := s.Query("SELECT * FROM users WHERE LOWER(name) = $1", strings.ToLower(name))
	if err != nil {
		return empty, errors.Join(ErrorInRequest, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return empty, ErrorUserNotFound
	}
	rows.Scan(&empty)

	return userFromQuery(rows)
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

func (s *postgresStore) initdb() error {
	pw, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), PasswordCost)
	if err != nil {
		return errors.Join(errors.New("failed to hash default password"), err)
	}

	sql := `
        CREATE TABLE IF NOT EXISTS users (
            id serial PRIMARY KEY,
            name varchar(24) NOT NULL UNIQUE,
            password varchar(128) NOT NULL,
            permissions integer NOT NULL default 0,
            createdAt Timestamp DEFAULT current_timestamp
        );

        DO $$
        BEGIN
            IF NOT EXISTS(SELECT 1 FROM users) THEN
                INSERT INTO users (id, name, password, permissions) VALUES (0, '%v', '%v', %v);
            END IF;
        END $$
    `
	sql = fmt.Sprintf(sql, defaultUsername, string(pw), models.PermissionsAdmin)
	_, err = s.Exec(sql)
	if err != nil {
		return errors.Join(errors.New("error when running initial db query"), err)
	}

	return nil
}

func userFromQuery(rows *sql.Rows) (models.User, error) {
	var user models.User
	err := rows.Scan(&user.Id, &user.Name, &user.Passwd, &user.Permissions, &user.CreatedAt)
	return user, err
}
