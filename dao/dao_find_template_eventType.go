package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "findsalon-backend/dbConfig"
    "findsalon-backend/dto"
)

func FindTemplateByEventType(eventType string) (dto.Template, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var template dto.Template
    err := dbConfig.DATABASE.Collection(dbConfig.TEMPLATES_COLLECTION).FindOne(ctx, bson.M{
        "EventType": eventType,
        "IsActive": true,
        "Deleted":  false,
    }).Decode(&template)
    if err != nil {
        return dto.Template{}, nil
    }
    return template, nil
}
