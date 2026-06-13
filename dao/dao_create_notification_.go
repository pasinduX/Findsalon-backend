package dao

import (
    "context"
    "time"

    "findsalon-backend/dbConfig"
    "findsalon-backend/dto"
)

func CreateNotification(n dto.Notification) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := dbConfig.DATABASE.Collection(dbConfig.NOTIFICATIONS_COLLECTION).InsertOne(ctx, n)
    return err
}
