package dto

import "time"

type Location struct {
	Latitude  float64 `json:"Latitude" bson:"Latitude"`
	Longitude float64 `json:"Longitude" bson:"Longitude"`
	Address   string  `json:"Address" bson:"Address"`
}

type GeoPoint struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

func NewGeoPoint(lat, lng float64) GeoPoint {
	return GeoPoint{Type: "Point", Coordinates: []float64{lng, lat}}
}

type Salon struct {
	SalonId       string    `json:"SalonId" bson:"SalonId"`
	OwnerId       string    `json:"OwnerId" bson:"OwnerId" validate:"required"`
	Name          string    `json:"Name" bson:"Name" validate:"required"`
	Description   string    `json:"Description" bson:"Description"`
	Address       string    `json:"Address" bson:"Address" validate:"required"`
	Area          string    `json:"Area" bson:"Area" validate:"required"`
	DistrictId    string    `json:"DistrictId" bson:"DistrictId" validate:"required"`
	DistrictName  string    `json:"DistrictName" bson:"DistrictName"`
	City          string    `json:"City" bson:"City" validate:"required"`
	CityId        string    `json:"CityId" bson:"CityId" validate:"required"`
	Phone         string    `json:"Phone" bson:"Phone"`
	Location      Location  `json:"Location" bson:"Location"`
	GeoLocation   GeoPoint  `json:"GeoLocation" bson:"GeoLocation"`
	CoverImageUrl string    `json:"CoverImageUrl" bson:"CoverImageUrl"`
	LogoUrl       string    `json:"LogoUrl" bson:"LogoUrl"`
	IsActive      bool      `json:"IsActive" bson:"IsActive"`
	CreatedAt     time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted       bool      `json:"Deleted" bson:"Deleted"`
}
