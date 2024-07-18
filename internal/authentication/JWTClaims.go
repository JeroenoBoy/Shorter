package authentication

import (
	"github.com/JeroenoBoy/Shorter/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserId             models.UserId `json:"userId"`
	HashedPasswordHash int           `json:"version"` // Its basically a verion right?
}
