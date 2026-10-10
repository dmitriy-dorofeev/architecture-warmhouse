package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func main() {
	http.HandleFunc("/temperature", temperatureByLocation)
	http.HandleFunc("/temperature/", temperatureByID)

	log.Println("temperature-api listening on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// GET /temperature?location=Living Room
func temperatureByLocation(w http.ResponseWriter, r *http.Request) {
	writeTemperature(w, r.URL.Query().Get("location"), "")
}

// GET /temperature/1
func temperatureByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/temperature/")
	if id == "" {
		http.NotFound(w, r)
		return
	}
	writeTemperature(w, "", id)
}

func writeTemperature(w http.ResponseWriter, location, sensorID string) {
	// If no location is provided, use a default based on sensor ID
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// If no sensor ID is provided, generate one based on location
	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	resp := TemperatureResponse{
		Value:       math.Round((15+rand.Float64()*15)*100) / 100,
		Unit:        "°C",
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Simulated temperature sensor",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
