package dbConfig

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(DATABASE_URL)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	DATABASE = client.Database(DATABASE_NAME)
	log.Printf("Connected to MongoDB database: %s", DATABASE_NAME)

	if err := ensureIndexes(context.Background()); err != nil {
		log.Printf("Warning: failed to create indexes: %v", err)
	}
}

func ensureIndexes(ctx context.Context) error {
	// 2dsphere index for salon geo queries
	salons := DATABASE.Collection(SALONS_COLLECTION)
	_, err := salons.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "GeoLocation", Value: "2dsphere"}},
		Options: options.Index().SetName("salon_geolocation_2dsphere").SetBackground(true),
	})
	if err != nil {
		log.Printf("Warning: could not create 2dsphere index: %v", err)
	}

	// Unique partial index on bookings (BarberId, StartTime) for direct-booking double-booking prevention
	bookings := DATABASE.Collection(BOOKINGS_COLLECTION)
	_, err = bookings.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "BarberId", Value: 1},
			{Key: "StartTime", Value: 1},
		},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.M{
				"Status":    bson.M{"$ne": "cancelled"},
				"StartTime": bson.M{"$gt": time.Time{}},
			}).
			SetName("unique_barber_starttime_active"),
	})
	if err != nil {
		log.Printf("Warning: could not create booking unique index: %v", err)
	}

	// Compound query index for booking window lookups
	_, err = bookings.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "BarberId", Value: 1},
			{Key: "StartTime", Value: 1},
			{Key: "EndTime", Value: 1},
		},
		Options: options.Index().SetName("barber_window_lookup"),
	})
	return err
}
