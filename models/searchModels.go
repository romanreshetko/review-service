package models

import "time"

type ReviewPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ReviewSort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type RangeInt struct {
	Min *int `json:"min"`
	Max *int `json:"max"`
}

type RangeFloat struct {
	Min *float64 `json:"min"`
	Max *float64 `json:"max"`
}

type ReviewFilters struct {
	CityID               *int64      `json:"city_id"`
	Season               *string     `json:"season"`
	Budget               *RangeInt   `json:"budget"`
	Tags                 *[]string   `json:"tags"`
	Rating               *RangeFloat `json:"rating"`
	WithKids             *bool       `json:"with_kids"`
	WithPets             *bool       `json:"with_pets"`
	ElderlyPeople        *bool       `json:"elderly_people"`
	LimitedMobility      *bool       `json:"limited_mobility"`
	PhysicallyChallenged *bool       `json:"physically_challenged"`
	TripType             *string     `json:"trip_type"`
	KeyWords             *string     `json:"key_words"`
}

type ReviewSearchRequest struct {
	Filters    ReviewFilters     `json:"filters"`
	Sort       *ReviewSort       `json:"sort"`
	Pagination *ReviewPagination `json:"pagination"`
}

type ReviewGeneralData struct {
	ID           int64     `json:"id"`
	AuthorID     int64     `json:"author_id"`
	CreationDate time.Time `json:"creation_date"`
	City         string    `json:"city"`
	MainPhoto    string    `json:"main_photo"`
	LikesNumber  int       `json:"likes_number"`
	ReviewMark   float64   `json:"review_mark"`
	TextStart    string    `json:"text_start"`
}
