package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zzzNitro/kenkkurrency/repository"
	"github.com/zzzNitro/kenkkurrency/usecase"
)

func main() {

	err := godotenv.Load("keychain.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Load the API key from environment variables
	weatherKey := os.Getenv("WEATHER_API_KEY")
	if weatherKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is not set")
	}

	// Initialize the weather repository with the API key
	weatherRepo := repository.NewApiRepository(weatherKey)

	// Create the usecase with the weather repository
	usecase := usecase.NewUsecase(weatherRepo)

	// Define routes for each concurrency demonstration
	http.HandleFunc("/api/weather/control", usecase.HandleControl)
	http.HandleFunc("/api/weather/waitgroup", usecase.HandleWaitGroup)
	http.HandleFunc("/api/weather/channels", usecase.HandleChannels)
	http.HandleFunc("/api/weather/mutex", usecase.HandleMutexes)

	// Start the HTTP server
	log.Println("Server starting on port :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
