package dao

import (
    "context"
    "time"

    "findsalon-backend/dbConfig"
    "findsalon-backend/dto"
)

func CreateDistrict(district dto.DistrictDTO) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := dbConfig.DATABASE.Collection(dbConfig.DISTRICTS_COLLECTION)
    _, err := collection.InsertOne(ctx, district)
    return err
}
