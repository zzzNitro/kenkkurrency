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

	newsKey := os.Getenv("CURRENTS_API_KEY")
	if newsKey == "" {
		log.Fatal("CURRENTS_API_KEY environment variable is not set")
	}

	weatherRepo := repository.NewWeatherRepository(weatherKey)
	newsRepo := repository.NewNewsRepository(newsKey)

	usecase := usecase.NewUsecase(weatherRepo, newsRepo)

	http.HandleFunc("/api/control", usecase.HandleControl)
	http.HandleFunc("/api/waitgroup", usecase.HandleWaitGroup)
	http.HandleFunc("/api/channels", usecase.HandleChannels)
	http.HandleFunc("/api/mutex", usecase.HandleMutexes)

	log.Println("Server starting on port :", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start: ", err)
	}

}
