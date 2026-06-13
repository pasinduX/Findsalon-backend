package dto

type RatingSummary struct {
	EntityId      string         `json:"EntityId"`
	EntityType    string         `json:"EntityType"`
	AverageRating float64        `json:"AverageRating"`
	TotalReviews  int            `json:"TotalReviews"`
	Distribution  map[string]int `json:"Distribution"`
	RecentReviews []Review       `json:"RecentReviews"`
}
