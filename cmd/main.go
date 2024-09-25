package main

import (
	"fmt"
	"net/http"

	"github.com/hiteshwadhwani/go-youtube-scrapper.git/internal/db"
	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/log"

	youtubeservice "github.com/hiteshwadhwani/go-youtube-scrapper.git/internal/youtube-service"
)

var logger = log.New()

func main() {
	config := &db.Config{
		Host:      "localhost",
		Port:      5432,
		User:      "postgres",
		Password:  "password",
		DbName:    "youtube-scrapper",
		TableName: "youtube_data",
	}

	db, err := db.New(config)
	if err != nil {
		logger.Error(fmt.Sprintf("Error connecting to database: %v", err))
	}

	client := &http.Client{}

	youtubeservice.RegisterHandlers(client, db, logger)

	http.ListenAndServe(":8080", nil)
}
