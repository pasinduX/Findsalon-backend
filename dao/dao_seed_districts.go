package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    "findsalon-backend/data"
    "findsalon-backend/dbConfig"
)

func SeedDistricts() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    collection := dbConfig.DATABASE.Collection(dbConfig.DISTRICTS_COLLECTION)

    for _, district := range data.Districts {
        _, err := collection.UpdateOne(
            ctx,
            bson.M{"id": district.ID},
            bson.M{"$set": district},
            options.Update().SetUpsert(true),
        )
        if err != nil {
            return err
        }
    }

    return nil
}
