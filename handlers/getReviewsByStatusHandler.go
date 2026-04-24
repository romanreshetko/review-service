package handlers

import (
	"encoding/json"
	"net/http"
	"review-service/models"
	"review-service/repository"
)

func (h *Handler) GetReviewsByStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "moderator" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	status := r.URL.Query().Get("status")

	reviews, err := repository.GetReviewsByStatus(h.db, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
