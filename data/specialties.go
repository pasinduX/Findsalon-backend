package data

import (
	_ "embed"
	"encoding/json"
	"log"

	"findsalon-backend/dto"
)

//go:embed specialties.json
var specialtiesJSON []byte

var Specialties []dto.Specialty

func init() {
	if err := json.Unmarshal(specialtiesJSON, &Specialties); err != nil {
		log.Fatalf("failed to parse specialties.json: %v", err)
	}
}
