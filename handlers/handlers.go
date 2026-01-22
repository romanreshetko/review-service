package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"review-service/models"
	"review-service/repository"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{db}
}
func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "user" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req models.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	sectionsJSON, err := json.Marshal(req.Sections)
	if err != nil {
		http.Error(w, "invalid json: sections", http.StatusBadRequest)
		return
	}
	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		http.Error(w, "invalid json: tags", http.StatusBadRequest)
		return
	}

	err = repository.CreateReview(h.db, req, claims.UserID, sectionsJSON, tagsJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
