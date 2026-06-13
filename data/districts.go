package data

import (
	_ "embed"
	"encoding/json"
	"log"

	"findsalon-backend/dto"
)

//go:embed districts.json
var districtsJSON []byte

var Districts []dto.DistrictDTO

func init() {
	if err := json.Unmarshal(districtsJSON, &Districts); err != nil {
		log.Fatalf("failed to parse districts.json: %v", err)
	}
}
