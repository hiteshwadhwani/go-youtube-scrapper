package youtubeservice

import (
	"database/sql"

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
	var args []interface{}
	query := `SELECT title, description, published_at, thumbnail_url, channel_title, created_at, updated_at FROM youtube_data`

	if params.search != "" {
		query += ` WHERE search_vector @@ to_tsquery('pg_catalog.english', $1::text) LIMIT $2::bigint OFFSET $3::bigint`
		args = append(args, params.search, params.limit, params.offset)
	} else {
		query += ` LIMIT $1::bigint OFFSET $2::bigint`
		args = append(args, params.limit, params.offset)
	}

	rows, err := r.db.Query(query, args...)

	if err != nil {
		r.logger.Error(err)
		return nil, err
	}

	defer rows.Close()

	response := make([]entity.YoutubeData, 0)

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
