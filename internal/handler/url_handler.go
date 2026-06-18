package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"url-shortner/internal/service"
)

type URLHandler struct {
	service service.URLService
}

func NewURLHandler(service service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func (h *URLHandler) Encode(w http.ResponseWriter, r *http.Request) {
	var req EncodeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %e", err))
		return
	}

	shortURL, err := h.service.Encode(r.Context(), req.URL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, EncodeResponse{
		ShortURL: shortURL,
	})
}

func (h *URLHandler) Decode(w http.ResponseWriter, r *http.Request) {
	var req DecodeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %e", err))
		return
	}

	url, err := h.service.Decode(r.Context(), req.ShortURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, DecodeResponse{
		URL: url,
	})
}
