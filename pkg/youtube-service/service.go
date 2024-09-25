package youtubeservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	entity "github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/entity"
)

type YoutubeService struct {
	apiKey     string
	searchKey  string
	maxResults int
	httpClient *http.Client
}

func New(httpClient *http.Client, apiKey string, searchKey string, maxResults int) *YoutubeService {
	return &YoutubeService{
		apiKey:     apiKey,
		searchKey:  searchKey,
		maxResults: maxResults,
		httpClient: httpClient,
	}
}

func (y *YoutubeService) GetSearchResults() ([]byte, error) {
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/search?part=snippet&maxResults=%d&q=%v&key=%v", y.maxResults, y.searchKey, y.apiKey)

	resp, err := y.httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (y *YoutubeService) GetVideoDetails(data []byte) []entity.YoutubeData {
	var youtubeData map[string]interface{}

	json.Unmarshal(data, &youtubeData)

	var entityItems []entity.YoutubeData

	for _, item := range youtubeData["items"].([]interface{}) {
		data := entity.YoutubeData{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		snippet, ok := itemMap["snippet"].(map[string]interface{})

		if !ok {
			continue
		}

		if title, ok := snippet["title"].(string); ok {
			data.Title = title
		}
		if description, ok := snippet["description"].(string); ok {
			data.Description = description
		}
		if publishedAt, ok := snippet["publishedAt"].(string); ok {
			data.PublishedAt = publishedAt
		}
		if channelTitle, ok := snippet["channelTitle"].(string); ok {
			data.ChannelTitle = channelTitle
		}

		if thumbnails, ok := snippet["thumbnails"].(map[string]interface{}); ok {
			for _, thumbnail := range thumbnails {
				if url, ok := thumbnail.(map[string]interface{})["url"].(string); ok {
					data.ThumbnailUrl = append(data.ThumbnailUrl, url)
				}
			}
		}

		entityItems = append(entityItems, data)
	}
	return entityItems
}

func (y *YoutubeService) ScheduleVideoDetailsUpdate(delay time.Duration, channel chan []entity.YoutubeData) {
	ticker := time.NewTicker(delay * time.Second)
	go func() {
		for range ticker.C {
			if data, err := y.GetSearchResults(); err == nil {
				channel <- y.GetVideoDetails(data)
			}
		}
	}()
}
