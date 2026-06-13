package dao

import (
	"findsalon-backend/dbConfig"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteUser(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.USERS_COLLECTION)
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"UserId": userId},
		bson.M{"$set": bson.M{
			"Deleted":   true,
			"IsActive":  false,
			"UpdatedAt": time.Now(),
		}},
	)
	return err
}
