package functions

import (
	"context"
	"log"
	"time"

	"findsalon-backend/data"
	"findsalon-backend/dbConfig"
	"findsalon-backend/dto"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// SeedSpecialties inserts the static specialty list into MongoDB only when the collection is empty.
func SeedSpecialties() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.SPECIALTIES_COLLECTION)

	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		log.Printf("Specialties already seeded (%d found), skipping", count)
		return nil
	}

	docs := make([]interface{}, len(data.Specialties))
	for i, s := range data.Specialties {
		docs[i] = s
	}
	if _, err := col.InsertMany(ctx, docs); err != nil {
		return err
	}
	log.Printf("Seeded %d specialties", len(data.Specialties))
	return nil
}

// SeedDistricts inserts the static district list into MongoDB only when the collection is empty.
func SeedDistricts() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.DISTRICTS_COLLECTION)

	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		log.Printf("Districts already seeded (%d found), skipping", count)
		return nil
	}

	docs := make([]interface{}, len(data.Districts))
	for i, d := range data.Districts {
		docs[i] = d
	}
	if _, err := col.InsertMany(ctx, docs); err != nil {
		return err
	}
	log.Printf("Seeded %d districts", len(data.Districts))
	return nil
}

// SeedDefaultTemplates ensures the default email templates exist.
func SeedDefaultTemplates() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	col := dbConfig.DATABASE.Collection(dbConfig.TEMPLATES_COLLECTION)
	defaults := []struct {
		eventType string
		name      string
		subject   string
		body      string
	}{
		{dto.EventBookingCreated, "Booking Confirmed", "Your booking is confirmed", DefaultBookingCreatedTemplate()},
		{dto.EventBookingCancelled, "Booking Cancelled", "Your booking has been cancelled", DefaultBookingCancelledTemplate()},
		{dto.EventBookingCompleted, "Visit Complete", "Thank you for your visit", DefaultBookingCompletedTemplate()},
	}
	for _, d := range defaults {
		count, _ := col.CountDocuments(ctx, bson.M{"EventType": d.eventType, "Deleted": false})
		if count > 0 {
			continue
		}
		tmpl := dto.Template{
			TemplateId:   uuid.New().String(),
			EventType:    d.eventType,
			Name:         d.name,
			Subject:      d.subject,
			BodyTemplate: d.body,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if _, err := col.InsertOne(ctx, tmpl); err != nil {
			log.Printf("SeedDefaultTemplates: failed to seed %s: %v", d.eventType, err)
		}
	}
	log.Println("Default email templates seeded")
}
