package usecase

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/zzzNitro/kenkkurrency/repository"
)

type Usecase struct {
	WeatherRepo repository.WeatherRepository
}

func NewUsecase(weatherRepo repository.WeatherRepository) *Usecase {
	return &Usecase{WeatherRepo: weatherRepo}
}

// HandleControl handles the API call sequentially with detailed validation
func (uc *Usecase) HandleControl(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Control request...")
	city := r.URL.Query().Get("city")

	// Step 1: Validate the city parameter
	if errorMessages := validateCitySequential(city); len(errorMessages) > 0 {
		log.Println("Validation Errors: ", errorMessages)
		http.Error(w, "Validation error: "+strings.Join(errorMessages, "; "), http.StatusBadRequest)
		return
	}

	// Step 2: Fetch weather data using the validated city
	weatherData, err := uc.WeatherRepo.GetWeatherByCity(city)
	if err != nil {
		http.Error(w, "Error fetching weather data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 3: Marshal and send the response
	response, err := json.Marshal(weatherData)
	if err != nil {
		http.Error(w, "Error marshalling response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// validateCitySequential performs detailed validation of the city string
func validateCitySequential(city string) []string {
	var errors []string
	trimmedCity := strings.TrimSpace(city)

	// Check if the city is empty
	if len(trimmedCity) == 0 {
		errors = append(errors, "City cannot be empty")
	}

	// Check for maximum length
	if len(trimmedCity) > 100 {
		errors = append(errors, "City name is too long")
	}

	// Check for valid characters
	for _, char := range trimmedCity {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			errors = append(errors, "City contains invalid characters")
			break
		}
	}

	return errors
}

func (uc *Usecase) HandleWaitGroup(w http.ResponseWriter, r *http.Request) {
	// Concurrency with WaitGroup
}

func (uc *Usecase) HandleChannels(w http.ResponseWriter, r *http.Request) {
	// Concurrency with Channels
}

func (uc *Usecase) HandleMutexes(w http.ResponseWriter, r *http.Request) {
	// Concurrency with Mutexes
}
