package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindBarberByEmail(email string) (dto.Barber, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var barber dto.Barber
	collection := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	err := collection.FindOne(ctx, bson.M{"Email": email, "Deleted": false}).Decode(&barber)
	return barber, err
}

func UpdateBarberUserId(barberId string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	_, err := collection.UpdateOne(ctx,
		bson.M{"BarberId": barberId},
		bson.M{"$set": bson.M{"UserId": userId}},
	)
	return err
}
