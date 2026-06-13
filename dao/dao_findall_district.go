package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "findsalon-backend/dbConfig"
    "findsalon-backend/dto"
)

func FindAllDistricts() ([]dto.DistrictDTO, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := dbConfig.DATABASE.Collection(dbConfig.DISTRICTS_COLLECTION)
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var districts []dto.DistrictDTO
    if err = cursor.All(ctx, &districts); err != nil {
        return nil, err
    }
    return districts, nil
}
