package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	TODO_PORT     string
	TODO_DBFILE   string
	WEB_DIR       string
	TODO_PASSWORD string
}

func New() (Config, error) {
	return Config{
		TODO_PORT:     getEnv("TODO_PORT", "7540"),
		TODO_DBFILE:   getExecutable("TODO_DBFILE"),
		WEB_DIR:       getEnv("WEB_DIR", "./web"),
		TODO_PASSWORD: os.Getenv("TODO_PASSWORD"),
	}, nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getExecutable(key string) string {
	dbFile := os.Getenv(key)

	if dbFile == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(fmt.Errorf("error getting current working directory: %w", err))
		}
		dbFile = filepath.Join(cwd, "scheduler.db")
	}

	return dbFile
}
