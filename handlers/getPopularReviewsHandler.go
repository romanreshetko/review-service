package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"review-service/repository"
	"time"
)

func (h *Handler) GetPopularReviewsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nocache := r.URL.Query().Get("nocache") == "true"
	ctx := r.Context()
	cacheKey := "popular_reviews"

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

	reviews, err := repository.GetPopularReviews(h.db)
	if err != nil {
		http.Error(w, "failed to find reviews", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(reviews)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	h.redis.Set(ctx, cacheKey, data, 20*time.Minute)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
