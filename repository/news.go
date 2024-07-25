package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type newsRepository struct {
	apiKey string
}

func NewNewsRepository(apiKey string) *newsRepository {
	return &newsRepository{apiKey: apiKey}
}

// GetNewsByCity fetches news data for a specific city
func (repo *newsRepository) GetNewsByCity(city string) ([]NewsData, error) {
	url := fmt.Sprintf("https://api.currentsapi.services/v1/latest-news?language=es&apiKey=%s", repo.apiKey)

	// Trim line jumps and spaces from the URL
	cleanURL := strings.TrimSpace(url)

	resp, err := http.Get(cleanURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse struct {
		Status string     `json:"status"`
		News   []NewsData `json:"news"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return apiResponse.News, nil
}

func (repo *newsRepository) GetNewsByCityWG(city string) ([]NewsData, error) {
	var wg sync.WaitGroup
	var result []NewsData
	var finalError error

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Println("Fetching news data for city: ", city)
		url := fmt.Sprintf("https://api.currentsapi.services/v1/latest-news?language=es&apiKey=%s", repo.apiKey)

		// Trim line jumps and spaces from the URL
		cleanURL := strings.TrimSpace(url)

		resp, err := http.Get(cleanURL)
		if err != nil {
			finalError = err
			return
		}
		defer resp.Body.Close()

		var apiResponse struct {
			Status string     `json:"status"`
			News   []NewsData `json:"news"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			finalError = err
			return
		}

		result = apiResponse.News
	}()

	wg.Wait()
	return result, finalError
}

func (repo *newsRepository) GetNewsByCityChan(city string) ([]NewsData, error) {
	resultChan := make(chan []NewsData)
	errorChan := make(chan error)

	go func() {
		url := fmt.Sprintf("https://api.currentsapi.services/v1/latest-news?language=es&apiKey=%s", repo.apiKey)

		// Trim line jumps and spaces from the URL
		cleanURL := strings.TrimSpace(url)

		resp, err := http.Get(cleanURL)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		var apiResponse struct {
			Status string     `json:"status"`
			News   []NewsData `json:"news"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			errorChan <- err
			return
		}

		resultChan <- apiResponse.News
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	}
}

func (repo *newsRepository) GetNewsByCityMx(city string) ([]NewsData, error) {
	return nil, nil
}
