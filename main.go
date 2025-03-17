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
