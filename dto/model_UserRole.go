package dto

import "time"

const (
	RoleAdmin     = "admin"
	Rolemoderator = "moderator"
	RoleUser      = "user"
	RoleSalonOwner = "salon_owner"
	RoleBarber    = "barber"
	RoleStaff     = "staff"
)

type UserRole struct {
	RoleId    string    `json:"RoleId" bson:"RoleId"`
	UserId    string    `json:"UserId" bson:"UserId" validate:"required"`
	SalonId   string    `json:"SalonId" bson:"SalonId"`
	Role      string    `json:"Role" bson:"Role" validate:"required"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted   bool      `json:"Deleted" bson:"Deleted"`
}
