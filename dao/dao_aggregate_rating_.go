package dao

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"
)

func AggregateRatingSummary(field string, value string) (dto.RatingSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := dbConfig.DATABASE.Collection(dbConfig.REVIEWS_COLLECTION)

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{field, value}, {"Deleted", false}, {"IsVisible", true}}}},
		{{"$group", bson.D{
			{"_id", nil},
			{"averageRating", bson.D{{"$avg", "$Rating"}}},
			{"totalReviews", bson.D{{"$sum", 1}}},
			{"count1", bson.D{{"$sum", bson.D{{"$cond", bson.A{bson.D{{"$eq", bson.A{"$Rating", 1}}}, 1, 0}}}}}},
			{"count2", bson.D{{"$sum", bson.D{{"$cond", bson.A{bson.D{{"$eq", bson.A{"$Rating", 2}}}, 1, 0}}}}}},
			{"count3", bson.D{{"$sum", bson.D{{"$cond", bson.A{bson.D{{"$eq", bson.A{"$Rating", 3}}}, 1, 0}}}}}},
			{"count4", bson.D{{"$sum", bson.D{{"$cond", bson.A{bson.D{{"$eq", bson.A{"$Rating", 4}}}, 1, 0}}}}}},
			{"count5", bson.D{{"$sum", bson.D{{"$cond", bson.A{bson.D{{"$eq", bson.A{"$Rating", 5}}}, 1, 0}}}}}},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return dto.RatingSummary{}, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return dto.RatingSummary{}, err
	}

	if len(results) == 0 {
		return dto.RatingSummary{}, nil
	}

	result := results[0]
	avg, _ := result["averageRating"].(float64)
	total, _ := result["totalReviews"].(int32)
	if total == 0 {
		if int64Value, ok := result["totalReviews"].(int64); ok {
			total = int32(int64Value)
		}
	}

	distribution := map[string]int{
		"1": toInt(result["count1"]),
		"2": toInt(result["count2"]),
		"3": toInt(result["count3"]),
		"4": toInt(result["count4"]),
		"5": toInt(result["count5"]),
	}

	recentReviews, _, err := FindAllReviews(bson.M{field: value, "IsVisible": true}, 0, 5)
	if err != nil {
		return dto.RatingSummary{}, err
	}

	return dto.RatingSummary{
		EntityId:      value,
		EntityType:    field,
		AverageRating: avg,
		TotalReviews:  int(total),
		Distribution:  distribution,
		RecentReviews: recentReviews,
	}, nil
}

func toInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}
