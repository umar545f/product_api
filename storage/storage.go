package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

func NewConnection(config *Config) (*gorm.DB, error) {
	// Create a connection string to the postgres server without specifying a database
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.SSLMode)

	// Connect to the server
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Create the database if it does not exist
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DBName))
	if err != nil && err.Error() != "pq: database \""+config.DBName+"\" already exists" {
		return nil, err
	}

	// Now connect to the database
	dsn = fmt.Sprintf(
		"host=%s port=%s password=%s user=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Password, config.User, config.DBName,
		config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil

}
