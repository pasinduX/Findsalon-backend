package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
)

func DeleteQuote(quoteId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.QUOTES_COLLECTION)
	update := bson.M{
		"$set": bson.M{
			"Deleted":   true,
			"UpdatedAt": time.Now(),
		},
	}
	_, err := collection.UpdateOne(ctx, bson.M{"QuoteId": quoteId}, update)
	return err
}
