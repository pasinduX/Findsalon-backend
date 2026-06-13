package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "findsalon-backend/dbConfig"
)

func DeleteTemplate(templateId string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := dbConfig.DATABASE.Collection(dbConfig.TEMPLATES_COLLECTION).UpdateOne(ctx, bson.M{
        "TemplateId": templateId,
        "Deleted":    false,
    }, bson.M{"$set": bson.M{"Deleted": true, "UpdatedAt": time.Now()}})
    return err
}
