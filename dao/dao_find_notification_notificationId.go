package dao

import (
    "context"
    "errors"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "findsalon-backend/dbConfig"
    "findsalon-backend/dto"
)

func FindNotificationById(notificationId string) (dto.Notification, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var notification dto.Notification
    err := dbConfig.DATABASE.Collection(dbConfig.NOTIFICATIONS_COLLECTION).FindOne(ctx, bson.M{
        "NotificationId": notificationId,
        "Deleted":        false,
    }).Decode(&notification)
    if err != nil {
        return dto.Notification{}, errors.New("notification not found")
    }
    return notification, nil
}
