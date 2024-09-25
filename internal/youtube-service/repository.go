package youtubeservice

import (
	"database/sql"
	"fmt"

	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/log"
)

type Repository interface {
	Get()
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

func (r *repository) Get() {
	rows, err := r.db.Query("SELECT * FROM youtube_data")

	if err != nil {
		r.logger.Error(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan()

		if err != nil {
			r.logger.Error("Error scanning data from database")
		}

	}

	cols, _ := rows.Columns()
	fmt.Print(cols)
}
