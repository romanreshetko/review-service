package handlers

import (
	"net/http"
	"review-service/models"
	"review-service/repository"
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
