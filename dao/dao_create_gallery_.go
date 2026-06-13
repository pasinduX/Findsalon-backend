package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func CreateGallery(gallery dto.Gallery) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.GALLERY_COLLECTION)
	_, err := collection.InsertOne(ctx, gallery)
	return err
}
