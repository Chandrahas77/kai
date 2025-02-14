# Golang Vulnerability Scanner

## Features
- **Golang Web Server** using `httprouter`
- **PostgreSQL Integration**
- **Dockerized Setup**
- **Database Migrations** using `goose`
- **Environment Variables Management**
- **Health Check for PostgreSQL**

## 📜 Prerequisites
Before setting up the project, ensure you have the following installed:

- **Go (v1.22 or later)** 
- **Docker** & **Docker Compose**

## Project Structure
kai-golang/
│── migrations/                # Goose migration files
│── internal/                   # Internal application logic
│   ├── database/               # Database connection and migrations
│   ├── handlers/               # HTTP route handlers
│   ├── config/                 # Configuration management
|   ├── services/               # business logic
│── .env                        # Environment variables (not committed)
│── Dockerfile                  # Docker build file
│── docker-compose.yml          # Docker Compose configuration
│── go.mod                      # Go modules dependencies
│── go.sum                      # Go modules checksum
│── main.go                     # Application entry point
│── README.md                   # Documentation
├── Makefile                    # Make file to run app, stop app,check logs, etc.,

---

## Installation & Setup

### Clone the Repository

git clone https://github.com/your-repo/kai-golang.git
cd kai

## Setup env variables 
for eg
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=vulnerabilities_db
DB_HOST=postgres
DB_PORT=5432
APP_PORT=8080

To run application use `make run`
To stop application use `make stop`
To build docker files use `make build`
To see docker logs use `make logs`
To manually enter PostgreSQLContainer use `make psql`

## Running migrations manually
Make sure you are root directory for running these commands

`goose -dir migrations postgres "{your_db_url} up"`

To rollback the previous migration use
`goose -dir migrations postgres "{your_db_url} down"`

## Ping check on terminal
`curl -v http://localhost:8080/ping`

## Unit tests
Unit tests are present only database level operations
Command to test database operations`go test ./internal/daos -cover`
