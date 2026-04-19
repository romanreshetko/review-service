package repository

import (
	"database/sql"
	"review-service/models"
)

func GetDraftReviews(db *sql.DB, user int64) ([]models.ReviewGeneralData, error) {
	rows, err := db.Query(`
		SELECT r.id, r.author_id, r.creation_date, c.city, r.main_photo, r.likes_number, r.review_mark, r.review_content->0->>'text' 
		FROM reviews r
		JOIN cities c ON r.city_id = c.id
		WHERE r.author_id = $1 AND r.status = 'draft'
`, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.ReviewGeneralData
	for rows.Next() {
		var review models.ReviewGeneralData
		if err := rows.Scan(&review.ID, &review.AuthorID, &review.CreationDate, &review.City, &review.MainPhoto, &review.LikesNumber, &review.ReviewMark, &review.TextStart); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
