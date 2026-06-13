package dto

type Specialty struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}
