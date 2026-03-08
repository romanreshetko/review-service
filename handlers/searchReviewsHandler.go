package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"review-service/models"
	"review-service/repository"
	"strconv"
)

func (h *Handler) SearchReviewsHandler(w http.ResponseWriter, r *http.Request) {

	var req models.ReviewSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	reviews, err := repository.SearchReviews(h.db, req)
	if err != nil {
		log.Println("search review error: ", err)
		http.Error(w, "search review error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func (h *Handler) GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	reviewID, err := strconv.ParseInt(r.URL.Query().Get("review_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect review_id", http.StatusBadRequest)
		return
	}
}
