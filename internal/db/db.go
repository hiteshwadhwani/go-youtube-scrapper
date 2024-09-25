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

func New(config *Config) (*sql.DB, error) {
	return createConnection(config)
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
	vectorSearchColumnName := "search_vector"

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
	tableCreationQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id SERIAL PRIMARY KEY, title TEXT, description TEXT, published_at TIMESTAMP, thumbnail_url TEXT[], channel_title TEXT, created_at TIMESTAMP, updated_at TIMESTAMP);", tableName)

	createFullTextSearch(db, tableName, vectorSearchColumnName)

	if _, err := db.Exec(tableCreationQuery); err != nil {
		return err
	}

	return nil
}

func createFullTextSearch(db *sql.DB, tableName string, vectorSearchColumnName string) error {
	if err := createFullTextSearchColumn(db, tableName, vectorSearchColumnName); err != nil {
		return err
	}

	if err := createFullTextSearchIndex(db, tableName, vectorSearchColumnName); err != nil {
		return err
	}

	if err := createFullTextSearchTrigger(db, tableName, vectorSearchColumnName); err != nil {
		return err
	}

	return nil
}

func createFullTextSearchColumn(db *sql.DB, tableName string, columnName string) error {
	_, err := db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN IF NOT EXISTS %s tsvector;", tableName, columnName))
	return err
}

func createFullTextSearchTrigger(db *sql.DB, tableName string, vectorSearchColumnName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE TRIGGER %s_fts_trigger BEFORE INSERT OR UPDATE ON %s FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger(%s, 'pg_catalog.english', title, description);", tableName, tableName, vectorSearchColumnName))
	return err
}

func createFullTextSearchIndex(db *sql.DB, tableName string, vectorSearchColumnName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s_fts_idx ON %s USING GIN(%s);", tableName, tableName, vectorSearchColumnName))
	return err
}
