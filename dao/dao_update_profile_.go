package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateProfile(userId string, updateData map[string]interface{}, authHeader string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.USERS_COLLECTION)
	updateData["UpdatedAt"] = time.Now().UTC()
	_, err := col.UpdateOne(ctx, bson.M{"UserId": userId, "Deleted": false}, bson.M{"$set": updateData})
	return err
}
