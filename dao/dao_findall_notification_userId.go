package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    "findsalon-backend/dbConfig"
    "findsalon-backend/dto"
)

func FindAllNotificationsByUserId(userId string, skip int64, limit int64) ([]dto.Notification, int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"UserId": userId, "Deleted": false}
    cursor, err := dbConfig.DATABASE.Collection(dbConfig.NOTIFICATIONS_COLLECTION).Find(ctx, filter, options.Find().SetSort(bson.M{"CreatedAt": -1}).SetSkip(skip).SetLimit(limit))
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(ctx)

    var notifications []dto.Notification
    if err := cursor.All(ctx, &notifications); err != nil {
        return nil, 0, err
    }

    total, err := dbConfig.DATABASE.Collection(dbConfig.NOTIFICATIONS_COLLECTION).CountDocuments(ctx, filter)
    if err != nil {
        return nil, 0, err
    }

    return notifications, int(total), nil
}
