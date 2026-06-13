package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"

	"go.mongodb.org/mongo-driver/bson"
)

func DeleteProfile(userId string, authHeader string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.USERS_COLLECTION)
	_, err := col.UpdateOne(ctx,
		bson.M{"UserId": userId},
		bson.M{"$set": bson.M{"Deleted": true, "UpdatedAt": time.Now().UTC()}},
	)
	return err
}
