package functions

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"

	"go.mongodb.org/mongo-driver/bson"
)

func UniqueCheck(collectionName, field, value string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(collectionName)
	count, err := collection.CountDocuments(ctx, bson.M{field: value, "Deleted": false})
	if err != nil {
		return false
	}
	return count == 0
}

func UniqueCheckRaw(collectionName, field, value string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := dbConfig.DATABASE.Collection(collectionName)
	count, err := collection.CountDocuments(ctx, bson.M{field: value})
	if err != nil {
		return false
	}
	return count == 0
}
