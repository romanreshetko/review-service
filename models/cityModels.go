package models

type City struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CityData struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Region        string  `json:"region"`
	ReviewsNumber int     `json:"reviews_number"`
	Mark          float64 `json:"mark"`
}

type CityName struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SearchClosestReviewRequest struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
