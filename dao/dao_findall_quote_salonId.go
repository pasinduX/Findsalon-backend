package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindAllQuotesBySalonId(salonId string) ([]dto.Quote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.QUOTES_COLLECTION)
	cursor, err := collection.Find(ctx, bson.M{"SalonId": salonId, "Deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var quotes []dto.Quote
	if err = cursor.All(ctx, &quotes); err != nil {
		return nil, err
	}
	return quotes, nil
}
