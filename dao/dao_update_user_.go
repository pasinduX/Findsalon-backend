package dao

import (
	"findsalon-backend/dbConfig"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateUser(userId string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update["UpdatedAt"] = time.Now()

	collection := dbConfig.DATABASE.Collection(dbConfig.USERS_COLLECTION)
	_, err := collection.UpdateOne(ctx, bson.M{"UserId": userId}, bson.M{"$set": update})
	return err
}
