package server

import (
	"fmt"
	"os"
	"strconv"
)

// Config ...
type Config struct {
	serverPort   string
	dbConnString string
	nWorkers     int
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

	nWorkersStr, exists := os.LookupEnv("NUMBER_OF_WORKERS")
	if !exists {
		return nil, fmt.Errorf("No NUMBER_OF_WORKERS in .env")
	}

	nWorkers, err := strconv.Atoi(nWorkersStr)
	if err != nil {
		return nil, err
	}

	return &Config{
		serverPort:   serverPort,
		dbConnString: dbConnString,
		nWorkers:     nWorkers,
	}, nil
}
