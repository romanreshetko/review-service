package handlers

import (
	"encoding/json"
	"net/http"
	"review-service/models"
	"review-service/repository"
	"strconv"
)

func (h *Handler) LikeReview(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "user" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	reviewID, err := strconv.ParseInt(r.URL.Query().Get("review_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect review_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		http.Error(w, "incorrect user_id", http.StatusUnauthorized)
		return
	}

	if r.Method == "POST" {
		err = repository.SaveLike(h.db, userID, reviewID)
		if err != nil {
			http.Error(w, "failed to save like", http.StatusInternalServerError)
			return
		}
	}
	if r.Method == "DELETE" {
		err = repository.DeleteLike(h.db, userID, reviewID)
		if err != nil {
			http.Error(w, "failed to delete like", http.StatusInternalServerError)
			return
		}
	}
	if r.Method == "GET" {
		isLiked, err := repository.GetLike(h.db, userID, reviewID)
		if err != nil {
			http.Error(w, "failed to find like", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]bool{"liked": isLiked})
		if err != nil {
			http.Error(w, "failed to find like", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetLikesByUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "user" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		http.Error(w, "incorrect user_id", http.StatusUnauthorized)
		return
	}
	reviews, err := repository.GetLikedReviews(h.db, userID)
	if err != nil {
		http.Error(w, "failed to find likes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}
