package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/JeroenoBoy/Shorter/internal/authentication"
	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/JeroenoBoy/Shorter/internal/web"
)

const ascii = `
  ___ _            _           
 / __| |_  ___ _ _| |_ ___ _ _ 
 \__ \ ' \/ _ \ '_|  _/ -_) '_|
 |___/_||_\___/_|  \__\___|_|  
`

func main() {
	var address = ":3000"

	println(ascii)
	println("Your self-hosted link shortner")
	println("")

	dataStore, err := datastore.NewPostgresStore(getPostgresConfig())
	if err != nil {
		panic(err)
	}

	jwtAuth := authentication.NewJWTAuthenticator(dataStore, getJWTSecret())

	handler := web.NewWebserver(jwtAuth, dataStore)
	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on port %v 🚀\n", address)
	println("Shorter is ready to roll!")
    println("")

	http.Serve(l, handler)
}

func getPostgresConfig() datastore.PostgresOptions {
	return datastore.PostgresOptions{
		Host:     envVar("SHORTER_POSTGRES_ADDRESS", "127.0.0.1"),
		Port:     envInt("SHORTER_POSTGRES_PORT", 5432),
		Database: envVar("SHORTER_POSTGRES_DATABASE", ""),
		UserName: envVar("SHORTER_POSTGRES_USER", ""),
		Password: envVar("SHORTER_POSTGRES_PASSWORD", ""),
		SSLMode:  envVar("SHORTER_POSTGRES_SSLMODE", "require"),
	}
}

func envBool(key string, def bool) bool {
	env, ok := os.LookupEnv(key)
	if ok {
		if env == "" {
			return false
		}
		b, err := strconv.ParseBool(env)
		if err != nil {
			panic(err)
		}
		return b
	} else {
		return def
	}
}

func envInt(key string, def int) int {
	env, ok := os.LookupEnv(key)
	if ok {
		i, err := strconv.Atoi(env)
		if err != nil {
			panic(err)
		}
		return i
	} else {
		return def
	}
}

func envVar(key string, def string) string {
	env, ok := os.LookupEnv(key)
	if ok {
		return env
	} else {
		return def
	}
}

func getJWTSecret() []byte {
	secret, ok := os.LookupEnv("SHORTER_JWT_SECRET")
	if !ok || len(secret) == 0 {
		fmt.Println("No $SHORTER_JWT_SECRET env variable provided! Generating temporary secret")
		available := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-=_+[]{}\\|;':\",.<>/?`~")
		chars := make([]rune, 120+rand.Intn(16))
		for i := range chars {
			chars[i] = available[rand.Intn(len(available)-1)]
		}
		return []byte(string(chars))
	}
	return []byte(secret)

}
