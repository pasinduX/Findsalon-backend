package dao

import (
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUserByUserId(userId string) (dto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user dto.User
	collection := dbConfig.DATABASE.Collection(dbConfig.USERS_COLLECTION)
	err := collection.FindOne(ctx, bson.M{"UserId": userId, "Deleted": false}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return dto.User{}, errors.New("user not found")
		}
		return dto.User{}, err
	}

	return user, nil
}
