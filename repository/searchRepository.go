package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"review-service/models"
	"strconv"
	"strings"
)

func SearchReviews(db *sql.DB, req models.ReviewSearchRequest) ([]models.ReviewGeneralData, error) {

	query, args := buildSearchQuery(req)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []models.ReviewGeneralData

	for rows.Next() {
		var r models.ReviewGeneralData
		err := rows.Scan(
			&r.ID,
			&r.AuthorID,
			&r.CreationDate,
			&r.City,
			&r.ReviewMark,
			&r.LikesNumber,
			&r.MainPhoto,
			&r.TextStart,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, nil
}

func buildSearchQuery(req models.ReviewSearchRequest) (string, []interface{}) {

	query := `
SELECT 
    r.id, r.author_id, r.creation_date, c.city, r.review_mark, r.likes_number,
    r.main_photo, r.review_content->0->>'text'
FROM reviews r 
JOIN cities c ON r.city_id = c.id
WHERE status = 'published'
`
	args := []interface{}{}
	idx := 1
	f := req.Filters

	if f.CityID != nil {
		query += " AND city_id = $" + strconv.Itoa(idx)
		args = append(args, *f.CityID)
		idx++
	}

	if f.Season != nil {
		query += " AND season = $" + strconv.Itoa(idx)
		args = append(args, *f.Season)
		idx++
	}

	if f.Budget != nil {
		if f.Budget.Min != nil {
			query += " AND budget >= $" + strconv.Itoa(idx)
			args = append(args, *f.Budget.Min)
			idx++
		}
		if f.Budget.Max != nil {
			query += " AND budget <= $" + strconv.Itoa(idx)
			args = append(args, *f.Budget.Max)
			idx++
		}
	}

	if f.Tags != nil && len(*f.Tags) > 0 {
		jsonTags, _ := json.Marshal(*f.Tags)
		query += " AND tags @> $" + strconv.Itoa(idx)
		args = append(args, jsonTags)
		idx++
	}

	if f.Rating != nil {
		if f.Rating.Min != nil {
			query += " AND review_mark >= $" + strconv.Itoa(idx)
			args = append(args, *f.Rating.Min)
			idx++
		}
		if f.Rating.Max != nil {
			query += " AND review_mark <= $" + strconv.Itoa(idx)
			args = append(args, *f.Rating.Max)
			idx++
		}
	}

	if f.WithKids != nil {
		query += " AND with_kids_flag = $" + strconv.Itoa(idx)
		args = append(args, *f.WithKids)
		idx++
	}

	if f.WithPets != nil {
		query += " AND with_pets_flag = $" + strconv.Itoa(idx)
		args = append(args, *f.WithPets)
		idx++
	}

	if f.ElderlyPeople != nil {
		query += " AND elderly_people_flag = $" + strconv.Itoa(idx)
		args = append(args, *f.ElderlyPeople)
		idx++
	}

	if f.LimitedMobility != nil {
		query += " AND limited_mobility_flag = $" + strconv.Itoa(idx)
		args = append(args, *f.LimitedMobility)
		idx++
	}

	if f.PhysicallyChallenged != nil {
		query += " AND physically_challenged_flag = $" + strconv.Itoa(idx)
		args = append(args, *f.PhysicallyChallenged)
		idx++
	}

	if f.TripType != nil {
		query += " AND trip_type = $" + strconv.Itoa(idx)
		args = append(args, *f.TripType)
		idx++
	}

	if f.KeyWords != nil {
		query += " AND review_tsv @@ plainto_tsquery('russian', $" + strconv.Itoa(idx) + ") AND review_tsv_flag = $" + strconv.Itoa(idx) + ")"
		args = append(args, *f.KeyWords)
		idx++
	}

	if req.Sort != nil {
		query += buildSortClause(*req.Sort)
	}
	if req.Pagination != nil {
		query += buildPaginationClause(*req.Pagination)
	}
	return query, args
}

func buildSortClause(sort models.ReviewSort) string {
	field := "creation_date"
	switch sort.Field {
	case "rating":
		field = "review_mark"
	case "date":
		field = "creation_date"
	case "popular":
		field = "likes_number"
	}

	dir := "DESC"
	if strings.ToUpper(sort.Direction) == "ASC" {
		dir = "ASC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", field, dir)
}

func buildPaginationClause(p models.ReviewPagination) string {
	limit := 20
	offset := 0

	if p.Limit > 0 && p.Limit < 100 {
		limit = p.Limit
	}
	if p.Offset > 0 {
		offset = p.Offset
	}

	return fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
}

func GetReviewByID(db *sql.DB, reviewID int64) (models.Review, error) {
	var r models.Review
	err := db.QueryRow(`
		SELECT 
    		id, author_id, creation_date, city_id, season, budget, tags,
    		transport_mark, cleanliness_mark, preservation_mark, safety_mark,
    		hospitality_mark, price_quality_ratio, review_mark,
    		with_kids_flag, with_pets_flag, pet, business_trip_flag,
    		physically_challenged_flag, limited_mobility_flag,
    		elderly_people_flag, special_diet_flag, likes_number,
    		trip_type, main_photo, status, review_content
		FROM reviews
WHERE id = $1
`, reviewID).Scan(
		&r.ID,
		&r.AuthorID,
		&r.CreationDate,
		&r.CityID,
		&r.Season,
		&r.Budget,
		&r.Tags,
		&r.TransportMark,
		&r.CleanlinessMark,
		&r.PreservationMark,
		&r.SafetyMark,
		&r.HospitalityMark,
		&r.PriceQualityRatio,
		&r.ReviewMark,
		&r.WithKidsFlag,
		&r.WithPetsFLag,
		&r.Pet,
		&r.BusinessTripFlag,
		&r.PhysicallyChallengedFlag,
		&r.LimitedMobilityFlag,
		&r.ElderlyPeopleFlag,
		&r.SpecialDietFlag,
		&r.LikesNumber,
		&r.MainPhoto,
		&r.Status,
		&r.TripType,
		&r.Sections,
	)
	if err != nil {
		return models.Review{}, err
	}

	return r, nil
}
