package config

import (
	"log"
	"os"
	"strconv"
)

// Config struct holds all configuration values
type Config struct {
	DBUrl      string
	MaxWorkers int
	GitHubRepo string
	Port       string
}

//Global instance of Config
var AppConfig Config

// LoadConfig loads environment variables and sets default values
func LoadConfig() {
	maxWorkers, err := strconv.Atoi(getEnv("MAX_WORKERS", "3"))
	if err != nil {
		log.Println("Invalid MAX_WORKERS value, using default (3)")
		maxWorkers = 3
	}

	AppConfig = Config{
		DBUrl:      getEnv("DATABASE_URL", "sqlite3://:memory:"),
		MaxWorkers: maxWorkers,
		GitHubRepo: getEnv("GITHUB_REPO", "https://github.com/velancio/vulnerability_scans"),
		Port:       getEnv("PORT", "8080"),
	}

	log.Printf("Config loaded: %+v", AppConfig)
}

// getEnv retrieves the environment variable or returns the default value
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
