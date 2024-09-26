## Youtube scrapper

API server built to continuously fetch and store the latest videos data from YouTube based on a predefined search query. The server provides endpoints for paginated access to the stored video data, as well as a search param for full text search


## Features

Fetches the latest videos every 10 seconds (customize it using env) using YouTube's Data API v3 and saves them into a PostgreSQL database.

API to return data reverse chronological order with pagination and full text search

indexes for quick retrieval
## API Reference

#### Get youtube data

```http
  GET /api/v1/youtube-data

  query params  - search
                - limit (default 10)
                - offset (default 0)

example (pagination) - /api/v1/youtube-data?limit=10&offset=1
example (full-text search) - /api/v1/youtube-data?search=tech
```
## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`DB_HOST`=localhost

`DB_PORT`=5432

`DB_USER`=postgres

`DB_PASSWORD`=password

`DB_NAME`=youtube-scrapper

`DB_TABLE_NAME`=youtube_data

`CRON_DELAY`=10 (10 seconds)

`YOUTUBE_API_KEY`= (https://console.cloud.google.com/apis)

`YOUTUBE_SEARCH_QUERY`=technology

`MAX_RESULTS`=5
## Run Locally

Clone the project

```bash
  git clone https://github.com/hiteshwadhwani/go-youtube-scrapper.git
```

Go to the project directory

```bash
  cd go-youtube-scrapper
```

Start container using docker compose 

```bash
  docker compose -f ./configs/docker-compose.dev.yaml up
```

Download dependencies

```bash
  go mod download
```

Run server

```bash
  go run cmd/main.go 
```

