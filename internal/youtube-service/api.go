package youtubeservice

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/log"
	"github.com/hiteshwadhwani/go-youtube-scrapper.git/pkg/types"
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
	http.HandleFunc("/api/v1/youtube-data", httpHandler.Get)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.Get(r)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := types.NewErrorResponse("Internal Server Error")
		json.NewEncoder(w).Encode(*response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := types.NewErrorResponse("Internal Server Error")
		json.NewEncoder(w).Encode(*response)
		return
	}

	response := types.NewSuccessResponse(*data)
	json.NewEncoder(w).Encode(*response)
}
