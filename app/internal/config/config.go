package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	TODO_PORT   string
	TODO_DBFILE string
}

func New() *Config {
	return &Config{
		TODO_PORT:   getEnv("TODO_PORT", "7540"),
		TODO_DBFILE: getExecutable("TODO_DBFILE"),
	}
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
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(fmt.Errorf("error getting appPath: %w", err))
		}
		dbFile = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	return dbFile
}
