package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host               string
	Port               int
	User               string
	Password           string
	DbName             string
	TableName          string
	CronDelay          int
	YoutubeApiKey      string
	YoutubeSearchQuery string
	MaxResults         int
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	DB_PORT, err := strconv.Atoi(os.Getenv("DB_PORT"))

	if err != nil {
		return nil, err
	}

	CRON_DELAY, err := strconv.Atoi(os.Getenv("CRON_DELAY"))

	if err != nil {
		return nil, err
	}

	MAX_RESULTS, err := strconv.Atoi(os.Getenv("MAX_RESULTS"))

	if err != nil {
		return nil, err
	}

	config := &Config{
		Host:               os.Getenv("DB_HOST"),
		Port:               DB_PORT,
		User:               os.Getenv("DB_USER"),
		Password:           os.Getenv("DB_PASSWORD"),
		DbName:             os.Getenv("DB_NAME"),
		TableName:          os.Getenv("DB_TABLE_NAME"),
		CronDelay:          CRON_DELAY,
		YoutubeApiKey:      os.Getenv("YOUTUBE_API_KEY"),
		YoutubeSearchQuery: os.Getenv("YOUTUBE_SEARCH_QUERY"),
		MaxResults:         MAX_RESULTS,
	}

	return config, nil
}
