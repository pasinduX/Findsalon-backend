package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
)

func DeleteUserRole(roleId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.USERROLES_COLLECTION)
	_, err := collection.UpdateOne(ctx, bson.M{"RoleId": roleId}, bson.M{"$set": bson.M{"Deleted": true, "UpdatedAt": time.Now()}})
	return err
}
