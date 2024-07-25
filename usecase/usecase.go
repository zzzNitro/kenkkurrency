package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/zzzNitro/kenkkurrency/repository"
)

type Usecase struct {
	WeatherRepo repository.WeatherRepository
	NewsRepo    repository.NewsRepository
}

func NewUsecase(
	weatherRepo repository.WeatherRepository,
	newsRepo repository.NewsRepository,
) *Usecase {
	return &Usecase{
		WeatherRepo: weatherRepo,
		NewsRepo:    newsRepo,
	}
}

// HandleControl handles the API call sequentially with detailed validation
func (uc *Usecase) HandleControl(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Control request...")
	startTime := time.Now()
	city := r.URL.Query().Get("city")

	// Step 1: Validate the city parameter
	if errorMessages := validateCity(city); len(errorMessages) > 0 {
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

	// Step 3: Fetch news data using the validated city
	newsData, err := uc.NewsRepo.GetNewsByCity(city)
	if err != nil {
		http.Error(w, "Error fetching news data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: Aggregate the responses
	response, err := aggregateResponses(weatherData, newsData, startTime)
	if err != nil {
		http.Error(w, "Error aggregating responses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (uc *Usecase) HandleWaitGroup(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling concurrent wait group request...")
	startTime := time.Now()
	city := r.URL.Query().Get("city")

	if errorMessages := validateCity(city); len(errorMessages) > 0 {
		log.Println("Validation Errors: ", errorMessages)
		http.Error(w, "Validation error: "+strings.Join(errorMessages, "; "), http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	var weatherData repository.WeatherData
	var newsData []repository.NewsData
	var errWeather, errNews error

	wg.Add(2) // Add two goroutines into the wait group

	go func() {
		defer wg.Done()
		weatherData, errWeather = uc.WeatherRepo.GetWeatherByCityWG(city)
	}()

	go func() {
		defer wg.Done()
		newsData, errNews = uc.NewsRepo.GetNewsByCityWG(city)
	}()

	wg.Wait() // Wait for both API calls to complete

	if errWeather != nil {
		http.Error(w, "Error fetching weather data: "+errWeather.Error(), http.StatusInternalServerError)
		return
	}
	if errNews != nil {
		http.Error(w, "Error fetching news data: "+errNews.Error(), http.StatusInternalServerError)
		return
	}

	response, err := aggregateResponses(weatherData, newsData, startTime)
	if err != nil {
		http.Error(w, "Error aggregating responses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (uc *Usecase) HandleChannels(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling Control request with channels...")
	startTime := time.Now()
	city := r.URL.Query().Get("city")

	// Step 1: Validate the city parameter
	if errorMessages := validateCity(city); len(errorMessages) > 0 {
		log.Println("Validation Errors: ", errorMessages)
		http.Error(w, "Validation error: "+strings.Join(errorMessages, "; "), http.StatusBadRequest)
		return
	}

	// Create channels for weather and news data
	weatherChan := make(chan repository.WeatherData)
	newsChan := make(chan []repository.NewsData)
	errorChan := make(chan error, 2) // Buffer to avoid blocking on error sends

	// Step 2: Fetch weather data using the validated city
	go func() {
		data, err := uc.WeatherRepo.GetWeatherByCityChan(city)
		if err != nil {
			errorChan <- err
			return
		}
		weatherChan <- data
	}()

	// Step 3: Fetch news data using the validated city
	go func() {
		data, err := uc.NewsRepo.GetNewsByCityChan(city)
		if err != nil {
			errorChan <- err
			return
		}
		newsChan <- data
	}()

	// Step 4: Wait for both operations to complete and check for errors
	var weatherData repository.WeatherData
	var newsData []repository.NewsData
	var errCount int

	for i := 0; i < 2; i++ { // We expect two responses
		select {
		case weather := <-weatherChan:
			weatherData = weather
		case news := <-newsChan:
			newsData = news
		case err := <-errorChan:
			log.Println("Error fetching data: ", err)
			errCount++
		}
	}

	if errCount > 0 {
		http.Error(w, "Failed to fetch data due to errors", http.StatusInternalServerError)
		return
	}

	// Step 5: Aggregate the responses
	response, err := aggregateResponses(weatherData, newsData, startTime)
	if err != nil {
		http.Error(w, "Error aggregating responses: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (uc *Usecase) HandleMutexes(w http.ResponseWriter, r *http.Request) {
	// Concurrency with Mutexes
}

// validateCity performs detailed validation of the city string
func validateCity(city string) []string {
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

func aggregateResponses(weatherData repository.WeatherData, newsData []repository.NewsData, startTime time.Time) ([]byte, error) {
	// Calculate elapsed time
	elapsed := time.Since(startTime)

	// Struct for the combined response
	type response struct {
		Weather repository.WeatherData `json:"weather"`
		News    []repository.NewsData  `json:"news"`
		Time    string                 `json:"time"`
	}

	resp := response{
		Weather: weatherData,
		News:    newsData,
		Time:    fmt.Sprintf("%.2fs", elapsed.Seconds()),
	}

	return json.Marshal(resp)
}
