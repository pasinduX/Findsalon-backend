package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func FindSalonBySalonId(salonId string) (dto.Salon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var salon dto.Salon
	collection := dbConfig.DATABASE.Collection(dbConfig.SALONS_COLLECTION)
	err := collection.FindOne(ctx, bson.M{"SalonId": salonId, "Deleted": false}).Decode(&salon)
	return salon, err
}
