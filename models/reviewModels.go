package models

import "time"

type AuthContext struct {
	UserID int64
	Role   string
}

type Place struct {
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type ReviewSection struct {
	Title  string   `json:"title"`
	Text   string   `json:"text"`
	Photos []string `json:"photos"`
	Places []Place  `json:"places"`
}

type CreateReviewRequest struct {
	CityID                   int             `json:"city_id"`
	Season                   string          `json:"season"`
	Budget                   int             `json:"budget"`
	Tags                     []string        `json:"tags"`
	TransportMark            int             `json:"transport_mark"`
	CleanlinessMark          int             `json:"cleanliness_mark"`
	PreservationMark         int             `json:"preservation_mark"`
	SafetyMark               int             `json:"safety_mark"`
	HospitalityMark          int             `json:"hospitality_mark"`
	PriceQualityRatio        int             `json:"price_quality_ratio"`
	WithKidsFlag             bool            `json:"with_little_kids_flag"`
	WithPetsFLag             bool            `json:"with_pets_flag"`
	Pet                      *string         `json:"pet"`
	PhysicallyChallengedFlag bool            `json:"physically_challenged_flag"`
	LimitedMobilityFlag      bool            `json:"limited_mobility_flag"`
	ElderlyPeopleFlag        bool            `json:"elderly_people_flag"`
	SpecialDietFlag          bool            `json:"special_diet_flag"`
	TripType                 string          `json:"type"`
	MainPhoto                string          `json:"main_photo"`
	Sections                 []ReviewSection `json:"sections"`
}

type Review struct {
	ID                       int             `json:"id"`
	AuthorID                 int             `json:"author_id"`
	CreationDate             time.Time       `json:"creation_date"`
	CityID                   int             `json:"city_id"`
	Season                   string          `json:"season"`
	Budget                   int             `json:"budget"`
	Tags                     []string        `json:"tags"`
	TransportMark            int             `json:"transport_mark"`
	CleanlinessMark          int             `json:"cleanliness_mark"`
	PreservationMark         int             `json:"preservation_mark"`
	SafetyMark               int             `json:"safety_mark"`
	HospitalityMark          int             `json:"hospitality_mark"`
	PriceQualityRatio        int             `json:"price_quality_ratio"`
	ReviewMark               float64         `json:"review_mark"`
	WithKidsFlag             bool            `json:"with_little_kids_flag"`
	WithPetsFLag             bool            `json:"with_pets_flag"`
	Pet                      string          `json:"pet"`
	PhysicallyChallengedFlag bool            `json:"physically_challenged_flag"`
	LimitedMobilityFlag      bool            `json:"limited_mobility_flag"`
	ElderlyPeopleFlag        bool            `json:"elderly_people_flag"`
	SpecialDietFlag          bool            `json:"special_diet_flag"`
	LikesNumber              int             `json:"likes_number"`
	MainPhoto                string          `json:"main_photo"`
	Status                   string          `json:"status"`
	TripType                 string          `json:"type"`
	Sections                 []ReviewSection `json:"sections"`
}
