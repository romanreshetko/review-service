package repository

import (
	"database/sql"
	"review-service/models"
)

func GetAllCities(db *sql.DB) ([]models.City, error) {
	rows, err := db.Query(`
		SELECT id, city, region, latitude, longitude 
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
			&city.Latitude,
			&city.Longitude,
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

func GetPopularCities(db *sql.DB) ([]models.CityName, error) {
	rows, err := db.Query(`
		SELECT id, city FROM cities
		ORDER BY reviews_number DESC
		LIMIT 10
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []models.CityName
	for rows.Next() {
		var city models.CityName
		err := rows.Scan(
			&city.ID,
			&city.Name,
		)
		if err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	return cities, nil
}
