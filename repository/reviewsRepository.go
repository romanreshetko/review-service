package repository

import (
	"database/sql"
	"errors"
	"review-service/models"
)

func CreateReview(db *sql.DB, req models.CreateReviewRequest, userId int64, sections, tags []byte) (err error) {
	reviewMark := float64(req.TransportMark+req.CleanlinessMark+req.PreservationMark+req.SafetyMark+req.HospitalityMark+req.PriceQualityRatio) / 6.0
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	_, err = tx.Exec(`INSERT INTO reviews (author_id, creation_date, city_id, season, budget, tags, 
                transport_mark, cleanliness_mark, preservation_mark, safety_mark, hospitality_mark, price_quality_ratio, review_mark, 
                with_kids_flag, with_pets_flag, pet, business_trip_flag, physically_challenged_flag, 
                limited_mobility_flag, elderly_people_flag, special_diet_flag, 
            	trip_type, main_photo, status, review_content) 
				VALUES ($1, NOW(), $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)`,
		userId, req.CityID, req.Season, req.Budget, tags,
		req.TransportMark, req.CleanlinessMark, req.PreservationMark, req.SafetyMark, req.HospitalityMark, req.PriceQualityRatio, reviewMark,
		req.WithKidsFlag, req.WithPetsFLag, SafeDeref(req.Pet), req.BusinessTripFlag, req.PhysicallyChallengedFlag,
		req.LimitedMobilityFlag, req.ElderlyPeopleFlag, req.SpecialDietFlag,
		req.TripType, req.MainPhoto, "published", sections)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE cities SET 
            mark = (mark * reviews_number + $1) / (reviews_number + 1),
            reviews_number = reviews_number + 1
            WHERE id = $2`,
		reviewMark, req.CityID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func DeleteReview(db *sql.DB, reviewID int64) (err error) {
	res, err := db.Exec(`DELETE FROM reviews WHERE id = $1`, reviewID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows == 1 {
		return nil
	}

	return err
}

func GetUserIdByReview(db *sql.DB, reviewID int64) (int64, error) {
	var id int64
	err := db.QueryRow(`SELECT author_id FROM reviews WHERE id = $1`, reviewID).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("incorrect reviewID")
		}
		return 0, err
	}

	return id, nil
}

func SafeDeref[T any](v *T) any {
	if v == nil {
		return nil
	}
	return *v
}
