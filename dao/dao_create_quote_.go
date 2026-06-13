package dao

import (
	"context"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func CreateQuote(quote dto.Quote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.QUOTES_COLLECTION)
	_, err := collection.InsertOne(ctx, quote)
	return err
}
