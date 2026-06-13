package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "findsalon-backend/dbConfig"
)

func CountUnreadNotifications(userId string) (int64, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    count, err := dbConfig.DATABASE.Collection(dbConfig.NOTIFICATIONS_COLLECTION).CountDocuments(ctx, bson.M{
        "UserId":  userId,
        "IsRead":  false,
        "Deleted": false,
    })
    if err != nil {
        return 0, err
    }
    return count, nil
}
