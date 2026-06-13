package dto

type CityDTO struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type DistrictDTO struct {
	ID     int       `json:"id" bson:"id"`
	Name   string    `json:"name" bson:"name"`
	Cities []CityDTO `json:"cities" bson:"cities"`
}
