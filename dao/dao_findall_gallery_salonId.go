package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindAllGalleryBySalonId(salonId string) ([]dto.Gallery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.GALLERY_COLLECTION)
	opts := options.Find().SetSort(bson.M{"SortOrder": 1})
	cursor, err := collection.Find(ctx, bson.M{"SalonId": salonId, "Deleted": false}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var gallery []dto.Gallery
	if err = cursor.All(ctx, &gallery); err != nil {
		return nil, err
	}
	return gallery, nil
}

func FindAllGalleryByBarberId(barberId string) ([]dto.Gallery, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.GALLERY_COLLECTION)
	opts := options.Find().SetSort(bson.M{"SortOrder": 1})
	cursor, err := collection.Find(ctx, bson.M{"BarberId": barberId, "Deleted": false}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var gallery []dto.Gallery
	if err = cursor.All(ctx, &gallery); err != nil {
		return nil, err
	}
	return gallery, nil
}
