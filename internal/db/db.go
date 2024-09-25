package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	Host      string
	Port      int
	User      string
	Password  string
	DbName    string
	TableName string
}

func createConnection(config *Config) (*sql.DB, error) {
	conn_string := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=disable", config.DbName, config.User, config.Password, config.Host, config.Port)

	db, err := sql.Open("postgres", conn_string)

	if err != nil {
		return nil, err
	}

	if err := ensureDb(db, config.DbName, config.TableName); err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func ensureDb(db *sql.DB, dbName string, tableName string) error {
	// dbExistsQuery := `SELECT EXISTS (SELECT datname FROM pg_catalog.pg_database WHERE datname = 'your_database_name');`

	dbExistsQuery := fmt.Sprintf("SELECT EXISTS (SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", dbName)

	var dbExists bool
	err := db.QueryRow(dbExistsQuery).Scan(&dbExists)

	if err != nil {
		return err
	}

	// create db if it does not exists
	if !dbExists {
		_, err = db.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			return err
		}
	}

	// create table if it does not exists
	tableCreationQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, title TEXT, description TEXT, published_at TIMESTAMP, thumbnail_url TEXT[], created_at TIMESTAMP, updated_at TIMESTAMP);", tableName)

	if _, err := db.Exec(tableCreationQuery); err != nil {
		return err
	}

	return nil
}

func New(config *Config) (*sql.DB, error) {
	return createConnection(config)
}
