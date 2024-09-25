package youtubeservice

import (
	"database/sql"
	"net/http"

	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/log"
)

type Handler interface {
	Get()
}

type handler struct {
	client  *http.Client
	logger  log.Logger
	service Service
}

func RegisterHandlers(client *http.Client, db *sql.DB, logger log.Logger) {
	repository := NewRepository(db, logger)
	service := NewService(repository)
	httpHandler := &handler{
		client:  client,
		logger:  logger,
		service: service,
	}
	http.HandleFunc("/api/v1", httpHandler.Get)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	h.service.Get()
	w.Write([]byte("Hello World"))
}
