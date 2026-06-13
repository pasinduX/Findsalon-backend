package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindProfileByEmail(email string, authHeader string) (dto.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.USERS_COLLECTION)
	var user dto.User
	err := col.FindOne(ctx, bson.M{"Email": email, "Deleted": false}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return dto.Profile{}, nil
	}
	if err != nil {
		return dto.Profile{}, err
	}
	return dto.Profile{
		UserId:          user.UserId,
		FullName:        user.FullName,
		Email:           user.Email,
		Phone:           user.Phone,
		AvatarUrl:       user.AvatarUrl,
		GoogleAvatarUrl: user.GoogleAvatarUrl,
		Provider:        user.Provider,
		Role:            user.Role,
		IsActive:        user.IsActive,
	}, nil
}
