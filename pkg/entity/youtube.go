package entity

import "time"

type YoutubeData struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	PublishedAt  string    `json:"published_at"`
	ThumbnailUrl []string  `json:"thumbnail_url"`
	ChannelTitle string    `json:"channel_title"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
