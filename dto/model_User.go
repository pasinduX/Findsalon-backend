package dto

import "time"

type User struct {
	UserId          string    `json:"UserId" bson:"UserId"`
	FullName        string    `json:"FullName" bson:"FullName" validate:"required"`
	Email           string    `json:"Email" bson:"Email" validate:"required,email"`
	Phone           string    `json:"Phone" bson:"Phone"`
	PhoneVerified   bool      `json:"PhoneVerified" bson:"PhoneVerified"`
	PhoneOTP        string    `json:"-" bson:"PhoneOTP"`
	PhoneOTPExpiry  time.Time `json:"-" bson:"PhoneOTPExpiry"`
	AvatarUrl       string    `json:"AvatarUrl" bson:"AvatarUrl"`
	GoogleAvatarUrl string    `json:"GoogleAvatarUrl" bson:"GoogleAvatarUrl"`
	Provider        string    `json:"Provider" bson:"Provider"`
	GoogleId        string    `json:"GoogleId" bson:"GoogleId"`
	Role            string    `json:"Role" bson:"Role"`
	IsActive        bool      `json:"IsActive" bson:"IsActive"`
	RefreshToken    string    `json:"RefreshToken" bson:"RefreshToken"`
	PasswordHash    string    `json:"PasswordHash" bson:"PasswordHash"`
	CreatedAt       time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt       time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
	Deleted         bool      `json:"Deleted" bson:"Deleted"`
}

type UserResponse struct {
	UserId          string    `json:"UserId"`
	FullName        string    `json:"FullName"`
	Email           string    `json:"Email"`
	Phone           string    `json:"Phone"`
	PhoneVerified   bool      `json:"PhoneVerified"`
	AvatarUrl       string    `json:"AvatarUrl"`
	GoogleAvatarUrl string    `json:"GoogleAvatarUrl"`
	Provider        string    `json:"Provider"`
	GoogleId        string    `json:"GoogleId"`
	Role            string    `json:"Role"`
	IsActive        bool      `json:"IsActive"`
	CreatedAt       time.Time `json:"CreatedAt"`
	UpdatedAt       time.Time `json:"UpdatedAt"`
}
