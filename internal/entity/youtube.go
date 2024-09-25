package entity

type YoutubeData struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	PublishedAt  string   `json:"published_at"`
	ThumbnailUrl []string `json:"thumbnail_url"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}
