// Package main implements a virtual pet simulation game where users can create,
// interact with, and care for digital pets.
//
// The virtual pet system simulates various pet attributes including:
//   - Health and hunger levels
//   - Mood and happiness
//   - Energy and activity status
//
// Users can perform various actions with their pets:
//   - Feed and care for their pet
//   - Play and interact with their pet
//   - Monitor their pet's status and well-being
//
// The simulation runs in real-time, with the pet's status changing based on time
// and user interactions.
package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/dblbee/govitualpet/ai"
	"github.com/dblbee/govitualpet/api"
)

var pet *ai.Pet

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the AI model and pet
	pet = ai.NewPet("Milo")

	// Initialize and start the API server
	apiServer := api.NewServer(pet)
	apiServer.Start(":3000")

	log.Println("Server started on port 3000")

	defer apiServer.Stop()
}
