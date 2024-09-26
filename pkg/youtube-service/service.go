package youtubeservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	apikeymanager "github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/api-key-manager"
	entity "github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/entity"
	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/log"
)

type YoutubeService struct {
	manager    *apikeymanager.Manager
	searchKey  string
	maxResults int
	httpClient *http.Client
	timeDelay  time.Duration
}

func New(httpClient *http.Client, manager *apikeymanager.Manager, searchKey string, maxResults int, timeDelay time.Duration) *YoutubeService {
	return &YoutubeService{
		manager:    manager,
		searchKey:  searchKey,
		maxResults: maxResults,
		httpClient: httpClient,
		timeDelay:  timeDelay,
	}
}

func (y *YoutubeService) GetSearchResults() ([]byte, error) {
	for {
		apiKey := y.manager.GetNextKey()
		if apiKey == "" {
			return nil, fmt.Errorf("api Key quota exceeded")
		}

		url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/search?part=snippet&maxResults=%d&q=%v&key=%v&publishedAfter=%v", y.maxResults, y.searchKey, apiKey, time.Now().UTC().Add(-1*y.timeDelay*time.Second).Format(time.RFC3339))

		resp, err := y.httpClient.Get(url)

		if resp.StatusCode == 403 {
			y.manager.MarkQuotaExceed()
			continue
		}

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
}

func (y *YoutubeService) GetVideoDetails(data []byte) []entity.YoutubeData {
	var youtubeData map[string]interface{}

	json.Unmarshal(data, &youtubeData)

	var entityItems []entity.YoutubeData

	_, ok := youtubeData["items"].([]interface{})

	if !ok {
		return entityItems
	}

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

func (y *YoutubeService) ScheduleVideoDetailsUpdate(delay time.Duration, channel chan []entity.YoutubeData, logger log.Logger) {
	ticker := time.NewTicker(delay * time.Second)
	go func() {
		for range ticker.C {
			if data, err := y.GetSearchResults(); err == nil {
				channel <- y.GetVideoDetails(data)
				logger.Info("Data fetched from youtube api")
			}
		}
	}()
}
