package dao

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    "findsalon-backend/dbConfig"
    "findsalon-backend/dto"
)

func FindAllTemplates() ([]dto.Template, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := dbConfig.DATABASE.Collection(dbConfig.TEMPLATES_COLLECTION).Find(ctx, bson.M{"Deleted": false}, options.Find().SetSort(bson.M{"CreatedAt": -1}))
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var templates []dto.Template
    if err := cursor.All(ctx, &templates); err != nil {
        return nil, err
    }
    return templates, nil
}
