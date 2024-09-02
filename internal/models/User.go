package models

import (
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type BcryptedPassword string
type UserId int

type User struct {
	Id          UserId           `json:"id"`
	Name        string           `json:"name"`
	Passwd      BcryptedPassword `json:"-"`
	Permissions Permissions      `json:"permissions"`
	CreatedAt   time.Time        `json:"created_at"`
}

func (p BcryptedPassword) Compare(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(p), []byte(pw))
}

func (u UserId) ToString() string {
	return strconv.Itoa(int(u))
}
