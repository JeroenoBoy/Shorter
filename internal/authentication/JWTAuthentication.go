package authentication

import (
	"context"
	"errors"
	"hash/fnv"
	"net/http"
	"time"

	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/JeroenoBoy/Shorter/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const tokenDuration = time.Minute

type JWTAuthentication struct {
	secret     []byte
	cookieName string
	store      datastore.Datastore
}

func NewJWTAuthenticator(store datastore.Datastore, secret []byte) *JWTAuthentication {
	return &JWTAuthentication{
		cookieName: "auth",
		secret:     secret,
		store:      store,
	}
}

func (j *JWTAuthentication) RetrieveClaims(token string) (JWTClaims, error) {
	var claims JWTClaims

	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method for JWT -- It must be HMAC")
		}

		return []byte(j.secret), nil
	})

	if err != nil {
		return JWTClaims{}, err
	}

	return claims, nil
}

func (j *JWTAuthentication) Sign(user models.User) (string, error) {
	h := fnv.New32()
	h.Write(user.Passwd)

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Shorter",
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserId:             user.Id,
		HashedPasswordHash: int(h.Sum32()),
	})

	return token.SignedString(j.secret)
}

func (j *JWTAuthentication) MiddilewareProvideUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer next.ServeHTTP(w, r) // Always go to the next one, this middileware should never catch errors

		cookie, err := r.Cookie(j.cookieName)
		if err != nil {
			return
		}

		token := cookie.Value
		claims, err := j.RetrieveClaims(token)
		if err != nil {
			return
		}
		r.WithContext(context.WithValue(r.Context(), CtxKeyClaims, claims))

		user, err := j.store.GetUser(claims.UserId)
		if err != nil {
			return
		}
		r.WithContext(context.WithValue(r.Context(), CtxKeyUser, user))
	})
}
