package http_transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/mizunaro/antifraud-service/internal/domain"
	"github.com/rs/zerolog/log"
)

type URLService interface {
	ProcessURL(ctx context.Context, rawURL string) (domain.URLCheck, error)
}

type Handler struct {
	service URLService
}

func NewHandler(service URLService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(r chi.Router) {
	r.Post("/api/v1/check", h.CheckURL)
}

type checkRequest struct {
	URL string `json:"url"`
}

type checkResponse struct {
	ID string `json:"id"`
}

func (h *Handler) CheckURL(w http.ResponseWriter, r *http.Request) {
	var req checkRequest

	// 1. Декодируем (используем NewDecoder для скорости)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// 2. Валидируем URL "на вшивость"
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		http.Error(w, "invalid url format", http.StatusBadRequest)
		return
	}

	// 3. Дергаем сервис
	check, err := h.service.ProcessURL(r.Context(), req.URL)
	if err != nil {
		log.Error().Err(err).Msg("failed to process url")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 4. Отдаем 202 Accepted
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(checkResponse{
		ID: check.ID.String(),
	})
}
