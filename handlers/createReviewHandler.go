package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	if err := r.ParseMultipartForm(64 << 20); err != nil {
		http.Error(w, "multipart too large", http.StatusBadRequest)
		return
	}

	reviewJSON := r.FormValue("review")
	if reviewJSON == "" {
		http.Error(w, "missing review", http.StatusBadRequest)
		return
	}

	var req models.CreateReviewRequest
	if err := json.Unmarshal([]byte(reviewJSON), &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	reviewID := uuid.New().String()
	baseDir := filepath.Join("uploads", "reviews", reviewID)
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		http.Error(w, "failed to create directory", http.StatusInternalServerError)
		return
	}

	for si, section := range req.Sections {
		for pi, photoID := range section.Photos {
			formKey := "photo_" + photoID
			file, header, err := r.FormFile(formKey)
			if err != nil {
				http.Error(w, "missing photo "+photoID, http.StatusBadRequest)
				return
			}

			func() {
				defer file.Close()

				ext := filepath.Ext(header.Filename)
				filename := photoID + uuid.New().String() + ext
				fullPath := filepath.Join(baseDir, filename)

				dst, err := os.Create(fullPath)
				if err != nil {
					http.Error(w, "cannot save photo: "+photoID, http.StatusInternalServerError)
					return
				}
				defer dst.Close()

				if _, err := io.Copy(dst, file); err != nil {
					http.Error(w, "cannot write photo: "+photoID, http.StatusInternalServerError)
					return
				}

				req.Sections[si].Photos[pi] = "/uploads/reviews/" + reviewID + "/" + filename
			}()
		}
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
