package dao

import (
    "context"
    "time"

    "findsalon-backend/dbConfig"
    "findsalon-backend/dto"
)

func CreateTemplate(t dto.Template) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := dbConfig.DATABASE.Collection(dbConfig.TEMPLATES_COLLECTION).InsertOne(ctx, t)
    return err
}
