package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
)

func UpdateUserRole(roleId string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.USERROLES_COLLECTION)
	update["UpdatedAt"] = time.Now()
	_, err := collection.UpdateOne(ctx, bson.M{"RoleId": roleId, "Deleted": false}, bson.M{"$set": update})
	return err
}
