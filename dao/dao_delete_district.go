package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "findsalon-backend/dbConfig"
)

func DeleteDistrict(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := dbConfig.DATABASE.Collection(dbConfig.DISTRICTS_COLLECTION)
    _, err := collection.DeleteOne(ctx, bson.M{"id": id})
    return err
}
