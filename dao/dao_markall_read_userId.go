package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "findsalon-backend/dbConfig"
)

func MarkAllNotificationsRead(userId string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := dbConfig.DATABASE.Collection(dbConfig.NOTIFICATIONS_COLLECTION).UpdateMany(ctx, bson.M{
        "UserId":  userId,
        "IsRead":  false,
        "Deleted": false,
    }, bson.M{"$set": bson.M{"IsRead": true, "UpdatedAt": time.Now()}})
    return err
}
