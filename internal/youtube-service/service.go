package youtubeservice

import (
	"net/http"

	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/constants"
	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/entity"
)

type Service interface {
	Get(r *http.Request) (*[]entity.YoutubeData, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

func (y *service) Get(r *http.Request) (*[]entity.YoutubeData, error) {
	search := r.URL.Query().Get("search")
	limit := r.URL.Query().Get("limit")

	if limit == "" {
		limit = constants.DEFAULT_PAGE_LIMIT
	}

	offset := r.URL.Query().Get("offset")

	if offset == "" {
		offset = constants.DEFAULT_PAGE_OFFSET
	}

	data, err := y.repository.Get(&GetRequestParams{
		search: search,
		limit:  limit,
		offset: offset,
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}
