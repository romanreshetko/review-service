package handlers

import (
	"encoding/json"
	"net/http"
	"review-service/models"
	"review-service/repository"
)

func (h *Handler) GetDraftsByUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "user" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	reviews, err := repository.GetDraftReviews(h.db, claims.UserID)
	if err != nil {
		http.Error(w, "failed to find drafts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
