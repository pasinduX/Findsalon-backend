package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindAllBarbersBySalonId(salonId string) ([]dto.Barber, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	cursor, err := collection.Find(ctx, bson.M{"SalonId": salonId, "Deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var barbers []dto.Barber
	if err = cursor.All(ctx, &barbers); err != nil {
		return nil, err
	}
	return barbers, nil
}
