package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

// FindAllBarbersByUserId returns every active Barber record linked to a given
// user across all salons. This powers the cross-salon membership view.
func FindAllBarbersByUserId(userId string) ([]dto.Barber, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.BARBERS_COLLECTION)
	cursor, err := collection.Find(ctx, bson.M{"UserId": userId, "Deleted": false})
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
