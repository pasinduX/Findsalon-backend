package dao

import (
    "context"
    "errors"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "findsalon-backend/dbConfig"
    "findsalon-backend/dto"
)

func FindTemplateById(templateId string) (dto.Template, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var template dto.Template
    err := dbConfig.DATABASE.Collection(dbConfig.TEMPLATES_COLLECTION).FindOne(ctx, bson.M{
        "TemplateId": templateId,
        "Deleted":    false,
    }).Decode(&template)
    if err != nil {
        return dto.Template{}, errors.New("template not found")
    }
    return template, nil
}
