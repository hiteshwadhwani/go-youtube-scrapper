package youtubeservice

import (
	"database/sql"
	"fmt"

	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/log"
	"github.com/lib/pq"

	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/entity"
)

type GetRequestParams struct {
	search string
	limit  string
	offset string
}

type Repository interface {
	Get(params *GetRequestParams) (*[]entity.YoutubeData, error)
}

type repository struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepository(db *sql.DB, logger log.Logger) *repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (r *repository) Get(params *GetRequestParams) (*[]entity.YoutubeData, error) {
	fmt.Println(params)
	rows, err := r.db.Query("SELECT title, description, published_at, thumbnail_url, channel_title, created_at, updated_at FROM youtube_data LIMIT $1 OFFSET $2", params.limit, params.offset)

	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	defer rows.Close()

	var response []entity.YoutubeData

	for rows.Next() {
		var data entity.YoutubeData
		err := rows.Scan(&data.Title, &data.Description, &data.PublishedAt, pq.Array(&data.ThumbnailUrl), &data.ChannelTitle, &data.CreatedAt, &data.UpdatedAt)

		if err != nil {
			r.logger.Error(err)
			continue
		}

		response = append(response, data)
	}

	return &response, nil
}
