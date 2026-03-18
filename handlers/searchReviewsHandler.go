package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"review-service/models"
	"review-service/repository"
	"strconv"
	"time"
)

func (h *Handler) SearchReviewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reviewID, err := strconv.ParseInt(r.URL.Query().Get("review_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect review_id", http.StatusBadRequest)
		return
	}

	nocache := r.URL.Query().Get("nocache") == "true"
	ctx := r.Context()
	cacheKey := "review" + strconv.FormatInt(reviewID, 10)

	if !nocache {
		cached, err := h.redis.Get(ctx, cacheKey).Result()

		if err == nil && cached != "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(cached)); err == nil {
				return
			}
			log.Printf("failed to write cache response: %v", err)
		}
		log.Println("No cache, going to DB")
	}

	review, err := repository.GetReviewByID(h.db, reviewID)
	if err != nil {
		http.Error(w, "review not found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(review)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	h.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
