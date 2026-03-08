package repository

import (
	"database/sql"
	"errors"
	"review-service/models"
)

func SaveLike(db *sql.DB, user, review int64) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	res, err := tx.Exec(`
		INSERT INTO review_likes (user_id, review_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
`, user, review)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 1 {
		_, err = tx.Exec(`
		UPDATE reviews
		SET likes_number = likes_number + 1
		WHERE id = $1
`, review)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func DeleteLike(db *sql.DB, user, review int64) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	res, err := tx.Exec(`
		DELETE FROM review_likes
		WHERE user_id = $1 AND review_id = $2
`, user, review)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 1 {
		_, err = tx.Exec(`
		UPDATE reviews
		SET likes_number = likes_number - 1
		WHERE id = $1
`, review)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func GetLike(db *sql.DB, user, review int64) (bool, error) {
	var id int64
	err := db.QueryRow(`
		SELECT user_id FROM review_likes
		WHERE user_id = $1 AND review_id = $2
`, user, review).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetLikedReviews(db *sql.DB, user int64) ([]models.ReviewGeneralData, error) {
	rows, err := db.Query(`
		SELECT r.id, r.author_id, r.creation_date, c.city, r.main_photo 
		FROM reviews r
		JOIN cities c ON r.city_id = c.id
		WHERE r.id IN 
			(SELECT review_id FROM review_likes
			WHERE user_id = $1)
`, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.ReviewGeneralData
	for rows.Next() {
		var review models.ReviewGeneralData
		if err := rows.Scan(&review.ID, &review.AuthorID, &review.CreationDate, &review.City, &review.MainPhoto); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
