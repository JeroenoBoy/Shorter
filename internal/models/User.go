package models

type Password []byte
type UserId int

type User struct {
	Id          UserId      `json:"id"`
	Name        string      `json:"name"`
	Passwd      Password    `json:"-"`
	Permissions Permissions `json:"permissions"`
}
