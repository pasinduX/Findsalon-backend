package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindBarberByBarberId(barberId string) (dto.Barber, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var barber dto.Barber
	collection := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	err := collection.FindOne(ctx, bson.M{"BarberId": barberId, "Deleted": false}).Decode(&barber)
	return barber, err
}
