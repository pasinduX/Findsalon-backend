package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindAllProfiles(authHeader string, skip int, limit int) ([]dto.Profile, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.USERS_COLLECTION)

	filter := bson.M{"Deleted": false}
	total64, err := col.CountDocuments(ctx, filter)
	if err != nil {
		return []dto.Profile{}, 0, err
	}

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "CreatedAt", Value: -1}})
	cursor, err := col.Find(ctx, filter, opts)
	if err != nil {
		return []dto.Profile{}, 0, err
	}
	defer cursor.Close(ctx)

	var profiles []dto.Profile
	for cursor.Next(ctx) {
		var user dto.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		profiles = append(profiles, dto.Profile{
			UserId:          user.UserId,
			FullName:        user.FullName,
			Email:           user.Email,
			Phone:           user.Phone,
			AvatarUrl:       user.AvatarUrl,
			GoogleAvatarUrl: user.GoogleAvatarUrl,
			Provider:        user.Provider,
			Role:            user.Role,
			IsActive:        user.IsActive,
		})
	}
	if profiles == nil {
		profiles = []dto.Profile{}
	}
	return profiles, int(total64), nil
}
