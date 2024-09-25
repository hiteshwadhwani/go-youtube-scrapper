package youtubeservice

import (
	"fmt"
	"io"
	"net/http"
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

func (y *YoutubeService) GetSearchResult() ([]byte, error) {
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
