package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hiteshwadhwani/go-youtube-scrapper.git/internal/config"
	"github.com/hiteshwadhwani/go-youtube-scrapper.git/internal/db"
	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/entity"
	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/log"
	"github.com/lib/pq"

	youtubeservice "github.com/hiteshwadhwani/go-youtube-scrapper.git/internal/youtube-service"

	youtubeCronService "github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/youtube-service"
)

var logger = log.New()

func main() {

	config, err := config.LoadConfig()

	if err != nil {
		logger.Error(fmt.Sprintf("Error loading config: %v", err))
		os.Exit(1)
	}

	dbConfig := &db.Config{
		Host:      config.Host,
		Port:      config.Port,
		User:      config.User,
		Password:  config.Password,
		DbName:    config.DbName,
		TableName: config.TableName,
	}

	db, err := db.New(dbConfig)
	if err != nil {
		logger.Error(fmt.Sprintf("Error connecting to database: %v", err))
		os.Exit(1)
	}

	defer db.Close()

	client := &http.Client{}

	data_channel := make(chan []entity.YoutubeData)
	// this go routine will schedule the cron job to fetch data from youtube api
	service := youtubeCronService.New(client, config.YoutubeApiKey, config.YoutubeSearchQuery, config.MaxResults, time.Duration(config.CronDelay))
	service.ScheduleVideoDetailsUpdate(time.Duration(config.CronDelay), data_channel, logger)
	// this go routine will listen to youtube api cron and insert data into database
	go insertYoutubeCronData(db, data_channel)

	// this will register the http handlers
	youtubeservice.RegisterHandlers(client, db, logger)

	http.ListenAndServe(":8080", nil)
}

func insertYoutubeCronData(db *sql.DB, data_channel chan []entity.YoutubeData) {
	for data := range data_channel {
		for _, youtubeData := range data {
			if _, err := db.Exec("INSERT INTO youtube_data (title, description, published_at, thumbnail_url, channel_title, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", youtubeData.Title, youtubeData.Description, youtubeData.PublishedAt, pq.Array(youtubeData.ThumbnailUrl), youtubeData.ChannelTitle, pq.FormatTimestamp(youtubeData.CreatedAt), pq.FormatTimestamp(youtubeData.UpdatedAt)); err != nil {
				logger.Error(fmt.Sprintf("Error inserting data into database: %v", err))
			}
		}
	}
}
