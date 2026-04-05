package repository

import (
	"database/sql"
	"review-service/models"
)

func GetAllCities(db *sql.DB) ([]models.City, error) {
	rows, err := db.Query(`
		SELECT id, city, region 
		FROM cities
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []models.City
	for rows.Next() {
		var city models.City
		err := rows.Scan(
			&city.ID,
			&city.Name,
			&city.Region,
		)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	return cities, nil
}

func GetCityByID(db *sql.DB, id int64) (models.CityData, error) {
	var city models.CityData
	err := db.QueryRow(`
		SELECT id, city, region, reviews_number, mark
		FROM cities
		WHERE id = $1
`, id).Scan(
		&city.ID,
		&city.Name,
		&city.Region,
		&city.ReviewsNumber,
		&city.Mark,
	)
	if err != nil {
		return models.CityData{}, err
	}

	return city, nil
}
