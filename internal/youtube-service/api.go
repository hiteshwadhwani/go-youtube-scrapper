package youtubeservice

import (
	"net/http"
)

type Handler interface {
	GetSearchResult()
}

type handler struct {
	service Service
}

func RegisterHandlers() {
	handler := &handler{}
	http.HandleFunc("/api/v1", handler.GetSearchResult)
}

func (h *handler) GetSearchResult(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))

}
