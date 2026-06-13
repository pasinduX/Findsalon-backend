package dto

import "time"

type Gallery struct {
	GalleryId string    `json:"GalleryId" bson:"GalleryId"`
	SalonId   string    `json:"SalonId" bson:"SalonId" validate:"required"`
	BarberId  *string   `json:"BarberId" bson:"BarberId"`
	ImageUrl  string    `json:"ImageUrl" bson:"ImageUrl" validate:"required"`
	Caption   string    `json:"Caption" bson:"Caption"`
	SortOrder int       `json:"SortOrder" bson:"SortOrder"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt"`
	Deleted   bool      `json:"Deleted" bson:"Deleted"`
}
