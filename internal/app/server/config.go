package server

import (
	"fmt"
	"os"
)

// Config ...
type Config struct {
	ServerPort      string
	DBConnString    string
}

// NewConfig - helper to init config
func NewConfig() (*Config, error) {
	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		return nil, fmt.Errorf("No SERVER_PORT in .env")
	}

	dbConnString, exists := os.LookupEnv("DB_CONN_STRING")
	if !exists {
		return nil, fmt.Errorf("No DB_CONN_STRING in .env")
	}

	return &Config{
		ServerPort:      serverPort,
		DBConnString:    dbConnString,
	}, nil
}
