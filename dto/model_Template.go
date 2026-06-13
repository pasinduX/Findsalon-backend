package dto

import "time"

type Template struct {
	TemplateId   string    `json:"TemplateId" bson:"TemplateId"`
	EventType    string    `json:"EventType" bson:"EventType" validate:"required"`
	Name         string    `json:"Name" bson:"Name" validate:"required"`
	Subject      string    `json:"Subject" bson:"Subject" validate:"required"`
	BodyTemplate string    `json:"BodyTemplate" bson:"BodyTemplate" validate:"required"`
	IsActive     bool      `json:"IsActive" bson:"IsActive"`
	CreatedAt    time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt    time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted      bool      `json:"Deleted" bson:"Deleted"`
}
