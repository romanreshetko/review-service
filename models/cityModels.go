package models

type City struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Region string `json:"region"`
}

type CityData struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Region        string  `json:"region"`
	ReviewsNumber int     `json:"reviews_number"`
	Mark          float64 `json:"mark"`
}
