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

	// client := &http.Client{}

	// youtubeService := youtubeservice.New(client, "AIzaSyBr06dXk45jp498vK4nO2CiTQD35u2cTjM", "car", 25)

	// data, _ := youtubeService.GetSearchResult()

	w.Write([]byte("Hello World"))

}
