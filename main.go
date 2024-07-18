package main

import (
	"os"
	"strconv"

	"github.com/JeroenoBoy/Shorter/internal/datastore"
	"github.com/JeroenoBoy/Shorter/internal/web"
)

func main() {

	dataStore, err := datastore.CreatePostgrsStore(getPostgresConfig())
	if err != nil {
		panic(err)
	}

	web.RunServer(dataStore)
}

func getPostgresConfig() datastore.PostgresOptions {
	return datastore.PostgresOptions{
		UserName: envVar("POSTGRES_USER", "root"),
		Password: envVar("POSTGRES_PASSWORD", ""),
		Ip:       envVar("POSTGRES_ADDRESS", "127.0.0.1"),
		Port:     envInt("POSTGRES_PORT", 5432),
		UseSSL:   envBool("POSTGRES_USESSL", false),
	}
}

func envBool(key string, def bool) bool {
	env, ok := os.LookupEnv(key)
	if ok {
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
