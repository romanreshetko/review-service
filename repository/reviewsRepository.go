package repository

import (
	"database/sql"
	"review-service/models"
)

func CreateReview(db *sql.DB, req models.CreateReviewRequest, userId, sections, tags []byte) error {
	_, err := db.Exec(`INSERT INTO reviews (author_id, creation_date, city_id, season, budget, tags, 
                transport_mark, cleanliness_mark, preservation_mark, safety_mark, hospitality_mark, price_quality_ratio, 
                with_kids_flag, with_pets_flag, pet, business_trip_flag, physically_challenged_flag, trip_type, main_photo, status, review_content) 
				VALUES ($1, NOW(), $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
		userId, req.CityID, req.Season, req.Budget, tags,
		req.TransportMark, req.CleanlinessMark, req.PreservationMark, req.SafetyMark, req.HospitalityMark, req.PriceQualityRatio,
		req.WithKidsFlag, req.WithPetsFLag, req.Pet, req.BusinessTripFlag, req.PhysicallyChallengedFlag, req.TripType,
		"", "published", sections)
	if err != nil {
		return err
	}
	//TODO Update city rating
	return nil
}
