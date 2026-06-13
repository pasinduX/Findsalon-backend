package dto

import "time"

type BlockType string

const (
	BlockTypeLunch    BlockType = "lunch"
	BlockTypeBreak    BlockType = "break"
	BlockTypeCritical BlockType = "critical"
)

type ScheduleBlock struct {
	BlockId   string    `json:"BlockId" bson:"BlockId"`
	BarberId  string    `json:"BarberId" bson:"BarberId" validate:"required"`
	SalonId   string    `json:"SalonId" bson:"SalonId" validate:"required"`
	StartTime time.Time `json:"StartTime" bson:"StartTime"`
	EndTime   time.Time `json:"EndTime" bson:"EndTime"`
	BlockType BlockType `json:"BlockType" bson:"BlockType" validate:"required"`
	Note      string    `json:"Note" bson:"Note"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
}

type CreateScheduleBlockRequest struct {
	BarberId  string    `json:"BarberId" validate:"required"`
	SalonId   string    `json:"SalonId" validate:"required"`
	Date      string    `json:"Date" validate:"required"`
	StartTime string    `json:"StartTime" validate:"required"`
	EndTime   string    `json:"EndTime" validate:"required"`
	BlockType BlockType `json:"BlockType" validate:"required"`
	Note      string    `json:"Note"`
}
