package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/JeroenoBoy/Shorter/api"
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

func (s *postgresStore) GetUsers(request models.PageRequest) (models.PaginatedUsers, error) {
	var empty models.PaginatedUsers
	if request.Page < 0 {
		return empty, ErrorInPreperation
	}
	if request.ItemsPerPage < 1 {
		return empty, ErrorInvalidItemCount
	}

	offset := request.Page * request.ItemsPerPage
	rows, err := s.Query(`
        SELECT u.*, total_users.total_users
        FROM (
            SELECT * FROM users OFFSET $1 LIMIT $2
        ) u
        CROSS JOIN (SELECT COUNT(*) AS total_users FROM users) total_users
    `, offset, request.ItemsPerPage)

	if err != nil {
		return empty, errors.Join(ErrorInRequest, err)
	}
	defer rows.Close()

	var paginatedUsers models.PaginatedUsers
	paginatedUsers.Users = make([]models.User, 0, request.ItemsPerPage)
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name, &user.Passwd, &user.Permissions, &user.CreatedAt, &paginatedUsers.ItemCount)
		if err != nil {
			return empty, errors.Join(ErrorInScan, err)
		}
		paginatedUsers.Users = append(paginatedUsers.Users, user)
	}

	paginatedUsers.PageCount = int(math.Ceil(float64(paginatedUsers.ItemCount) / float64(request.ItemsPerPage)))
	paginatedUsers.CurrentPage = request.Page

	return paginatedUsers, nil
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

func (s *postgresStore) GetLink(id models.LinkId) (models.ShortLink, error) {
	var empty models.ShortLink
	rows, err := s.Query("SELECT * FROM links WHERE id = $1", id)
	if err != nil {
		return empty, errors.Join(ErrorInRequest, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return empty, ErrorLinkNotFound
	}

	return linkFromQuery(rows)
}

func (s *postgresStore) GetUserLinks(userId models.UserId) ([]models.ShortLink, error) {
	rows, err := s.Query("SELECT * FROM links WHERE owner_id = $1", userId)
	if err != nil {
		return nil, errors.Join(ErrorInRequest, err)
	}
	defer rows.Close()

	result := make([]models.ShortLink, 0, 32)
	for rows.Next() {
		link, err := linkFromQuery(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, link)
	}

	return result, nil
}

func (s *postgresStore) GetAllLinks(request models.PageRequest) (models.PaginatedLinks, error) {
	var empty models.PaginatedLinks
	if request.Page < 1 {
		return empty, ErrorInPreperation
	}
	if request.ItemsPerPage < 1 {
		return empty, ErrorInvalidItemCount
	}

	offset := (request.Page - 1) * request.ItemsPerPage
	rows, err := s.Query(`
        SELECT u.*, total_links.total_links
        FROM (
            SELECT * FROM links OFFSET $1 LIMIT $2
        ) u
        CROSS JOIN (SELECT COUNT(*) AS total_links FROM links) total_links
    `, offset, request.ItemsPerPage)

	if err != nil {
		return empty, errors.Join(ErrorInRequest, err)
	}
	defer rows.Close()

	var paginatedLinks models.PaginatedLinks
	paginatedLinks.Links = make([]models.ShortLink, 0, request.ItemsPerPage)
	for rows.Next() {
		var link models.ShortLink
		err := rows.Scan(&link.Id, &link.Link, &link.Owner, &link.Target, &link.Redirects, &link.CreatedAt, &link.LastUsed, &paginatedLinks.ItemCount)
		if err != nil {
			return empty, errors.Join(ErrorInScan, err)
		}
		paginatedLinks.Links = append(paginatedLinks.Links, link)
	}

	paginatedLinks.PageCount = int(math.Ceil(float64(paginatedLinks.ItemCount) / float64(request.ItemsPerPage)))
	paginatedLinks.CurrentPage = request.Page

	return paginatedLinks, nil
}

func (s *postgresStore) CreateLink(owner models.UserId, link *string, target string) (models.ShortLink, error) {
	var empty models.ShortLink
	var err error
	var rows *sql.Rows

	if link != nil && len(*link) > 0 {
		rows, err = s.Query("INSERT INTO links (owner_id, link, target) VALUES ($1, $2, $3) RETURNING *", owner, *link, target)
	} else {
		rows, err = s.Query("INSERT INTO links (owner_id, target) VALUES ($1, $2) RETURNING *", owner, target)
	}

	if err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key") {
			err = errors.Join(ErrorDuplicateKey, err)
		}
		return empty, errors.Join(ErrorInRequest, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return empty, ErrorLinkNotFound
	}
	return linkFromQuery(rows)
}

func (s *postgresStore) UpdateLink(id models.LinkId, updateReq api.UpdateLinkRequest) (models.ShortLink, error) {
	rows, err := s.Query("UPDATE links SET link = $1, target = $2 WHERE id = $3 RETURNING *", updateReq.Link, updateReq.Target, id)
	if err != nil {
		if strings.Contains(err.Error(), "pq: duplicate key") {
			err = errors.Join(ErrorDuplicateKey, err)
		}
		return models.ShortLink{}, err
	}

	if !rows.Next() {
		return models.ShortLink{}, ErrorLinkNotFound
	}

	return linkFromQuery(rows)
}

func (s *postgresStore) DeleteLink(id models.LinkId) error {
	result, err := s.Exec("DELETE FROM links WHERE id = $1", id)
	if err != nil {
		return errors.Join(ErrorInRequest, err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrorLinkNotFound
	}
	return nil
}

func (s *postgresStore) GetLinkTargetAndIncreaseRedirects(link string) (string, error) {
	rows, err := s.Query("UPDATE links SET redirects = redirects + 1, last_used = now() WHERE link = $1 RETURNING target", link)
	if err != nil {
		return "", errors.Join(ErrorInRequest, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return "", ErrorLinkNotFound
	}

	var target string
	err = rows.Scan(&target)
	if err != nil {
		return "", err
	}
	return target, nil
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
            created_at Timestamp NOT NULL DEFAULT current_timestamp
        );

        CREATE TABLE IF NOT EXISTS links (
            id serial PRIMARY KEY,
            link varchar(24) UNIQUE NOT NULL DEFAULT substr(md5(random()::text), 1, 8),
            owner_id int NOT NULL,
            target varchar(512) NOT NULL,
            redirects int NOT NULL DEFAULT 0,
            created_at Timestamp NOT NULL DEFAULT current_timestamp,
            last_used Timestamp DEFAULT null,
            CONSTRAINT fk_owner_id FOREIGN KEY(owner_id) REFERENCES users(id)
        );
        
        CREATE INDEX IF NOT EXISTS idx_links_link ON links (link);

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
	if err != nil {
		return user, errors.Join(ErrorInScan, err)
	} else {
		return user, nil
	}
}

func linkFromQuery(rows *sql.Rows) (models.ShortLink, error) {
	var link models.ShortLink
	err := rows.Scan(&link.Id, &link.Link, &link.Owner, &link.Target, &link.Redirects, &link.CreatedAt, &link.LastUsed)
	if err != nil {
		return link, errors.Join(ErrorInScan, err)
	} else {
		return link, nil
	}
}
