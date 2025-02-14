package config

import (
	"fmt"
	"os"
)

// DBConfig holds database configuration details
type DBConfig struct {
	User       string
	Password   string
	Name       string
	Host       string
	Port       string
	ServerPort string
}

// LoadConfig reads environment variables and returns a DBConfig instance
func LoadConfig() (*DBConfig, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	serverPort := os.Getenv("APP_PORT")

	if dbUser == "" || dbPassword == "" || dbName == "" || dbHost == "" || dbPort == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return &DBConfig{
		User:       dbUser,
		Password:   dbPassword,
		Name:       dbName,
		Host:       dbHost,
		Port:       dbPort,
		ServerPort: serverPort,
	}, nil
}

// GetDSN constructs the database connection string
func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.Name)
}
