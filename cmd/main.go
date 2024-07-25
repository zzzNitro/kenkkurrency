package main

import (
	"log"
	"net/http"

	"github.com/zzzNitro/kenkkurrency/repository"
	"github.com/zzzNitro/kenkkurrency/usecase"
)

func main() {
	repo := repository.NewRepository()
	uc := usecase.NewUseCase(repo)

	// Set up routes for each concurrency demonstration
	http.HandleFunc("/api/control", uc.HandleControl)
	http.HandleFunc("/api/waitgroup", uc.HandleWaitGroup)
	http.HandleFunc("/api/channels", uc.HandleChannels)
	http.HandleFunc("/api/mutexes", uc.HandleMutexes)

	// Server configuration
	log.Println("Server starting on port :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
