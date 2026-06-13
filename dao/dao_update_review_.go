package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"findsalon-backend/dbConfig"
)

func UpdateReview(reviewId string, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update["UpdatedAt"] = time.Now()
	updateDoc := bson.M{"$set": update}

	collection := dbConfig.DATABASE.Collection(dbConfig.REVIEWS_COLLECTION)
	_, err := collection.UpdateOne(ctx, bson.M{"ReviewId": reviewId, "Deleted": false}, updateDoc)
	return err
}
