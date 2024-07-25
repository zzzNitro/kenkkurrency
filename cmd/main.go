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
		log.Println("No .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	weatherKey := os.Getenv("WEATHER_API_KEY")
	if weatherKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is not set")
	}

	weatherRepo := repository.NewApiRepository(weatherKey)

	usecase := usecase.NewUsecase(weatherRepo)

	http.HandleFunc("/api/weather/control", usecase.HandleControl)
	http.HandleFunc("/api/weather/waitgroup", usecase.HandleWaitGroup)
	http.HandleFunc("/api/weather/channels", usecase.HandleChannels)
	http.HandleFunc("/api/weather/mutex", usecase.HandleMutexes)

	log.Println("Server starting on port :", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start: ", err)
	}

}
