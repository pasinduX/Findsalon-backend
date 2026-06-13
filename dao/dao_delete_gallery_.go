package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
)

func DeleteGallery(galleryId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.GALLERY_COLLECTION)
	update := bson.M{
		"$set": bson.M{
			"Deleted": true,
		},
	}
	_, err := collection.UpdateOne(ctx, bson.M{"GalleryId": galleryId}, update)
	return err
}
