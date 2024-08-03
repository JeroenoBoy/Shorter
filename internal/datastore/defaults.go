package datastore

import "golang.org/x/crypto/bcrypt"

const (
	defaultUsername = "admin"
	defaultPassword = "Shorter"
	PasswordCost    = bcrypt.DefaultCost
)
