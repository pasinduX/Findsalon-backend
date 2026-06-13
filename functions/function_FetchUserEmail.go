package functions

import (
	"context"
	"log"
	"time"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"go.mongodb.org/mongo-driver/bson"
)

// FetchUserEmailAndName returns email and full name for a userId.
// Replaces the HTTP call to auth-service in the old notification-ms.
func FetchUserEmailAndName(userId string) (email, fullName string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.USERS_COLLECTION)
	var user dto.User
	if err = col.FindOne(ctx, bson.M{"UserId": userId, "Deleted": false}).Decode(&user); err != nil {
		log.Printf("FetchUserEmailAndName: user %s not found: %v", userId, err)
		return "", "", err
	}
	return user.Email, user.FullName, nil
}
