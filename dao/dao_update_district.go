package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "findsalon-backend/dbConfig"
)

func UpdateDistrict(id int, update bson.M) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := dbConfig.DATABASE.Collection(dbConfig.DISTRICTS_COLLECTION)
    _, err := collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": update})
    return err
}
