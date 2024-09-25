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
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DbName:   "youtube-scrapper",
	}

	db, err := db.New(config)
	if err != nil {
		logger.Error(fmt.Sprintf("Error connecting to database: %v", err))
	}

	fmt.Print(db)

	youtubeservice.RegisterHandlers()

	http.ListenAndServe(":8080", nil)
}
