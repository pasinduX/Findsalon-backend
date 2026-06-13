package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"findsalon-backend/dbConfig"
)

func DeleteReview(reviewId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"Deleted": true, "UpdatedAt": time.Now()}}
	collection := dbConfig.DATABASE.Collection(dbConfig.REVIEWS_COLLECTION)
	_, err := collection.UpdateOne(ctx, bson.M{"ReviewId": reviewId, "Deleted": false}, update)
	return err
}
