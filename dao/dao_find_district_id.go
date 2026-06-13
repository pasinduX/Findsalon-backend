package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "findsalon-backend/dbConfig"
    "findsalon-backend/dto"
)

func FindDistrictById(id int) (dto.DistrictDTO, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var district dto.DistrictDTO
    collection := dbConfig.DATABASE.Collection(dbConfig.DISTRICTS_COLLECTION)
    err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&district)
    return district, err
}
