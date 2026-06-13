package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "findsalon-backend/dbConfig"
)

func UpdateNotification(notificationId string, update bson.M) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    update["UpdatedAt"] = time.Now()
    _, err := dbConfig.DATABASE.Collection(dbConfig.NOTIFICATIONS_COLLECTION).UpdateOne(ctx, bson.M{
        "NotificationId": notificationId,
        "Deleted":        false,
    }, bson.M{"$set": update})
    return err
}
