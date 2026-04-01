package handlers

import (
	"net/http"
	"review-service/models"
	"review-service/repository"
	"slices"
	"strconv"
)

func (h *Handler) DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	reviewID, err := strconv.ParseInt(r.URL.Query().Get("review_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect review_id", http.StatusBadRequest)
		return
	}

	userID, err := repository.GetUserIdByReview(h.db, reviewID)
	if err != nil {
		if err.Error() == "incorrect reviewID" {
			http.Error(w, "review not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error getting review", http.StatusInternalServerError)
		return
	}

	if claims.Role != "user" || claims.UserID != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err = repository.DeleteReview(h.db, reviewID)
	if err != nil {
		http.Error(w, "error deleting review", http.StatusInternalServerError)
		return
	}

	//TODO comments delete
}

func (h *Handler) UpdateReviewStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims, ok := r.Context().Value("claims").(models.AuthContext)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	status := r.URL.Query().Get("status")
	validStatuses := []string{"published", "moderating", "blocked", "draft", "reported", "blocked_reported", "undefined", "moderation_error"}
	if !slices.Contains(validStatuses, status) {
		http.Error(w, "incorrect status", http.StatusBadRequest)
		return
	}

	reviewID, err := strconv.ParseInt(r.URL.Query().Get("review_id"), 10, 64)
	if err != nil {
		http.Error(w, "incorrect review_id", http.StatusBadRequest)
		return
	}

	if claims.Role == "user" && status != "reported" && status != "blocked_reported" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if claims.Role != "user" && claims.Role != "moderator" && claims.Role != "admin" && claims.Role != "service" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	err = repository.UpdateReviewStatus(h.db, reviewID, status)
	if err != nil {
		if err.Error() == "review not found" {
			http.Error(w, "review not found", http.StatusNotFound)
			return
		}
		http.Error(w, "error updating review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
