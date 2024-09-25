package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func newDb(config *Config) (*sql.DB, error) {
	conn_string := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=disable", config.DbName, config.User, config.Password, config.Host, config.Port)

	db, err := sql.Open("postgres", conn_string)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func New(config *Config) (*sql.DB, error) {
	return newDb(config)
}
