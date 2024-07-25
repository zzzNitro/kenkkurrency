package repository

/*
Example of weather response from the API:
{
    "location": {
        "name": "Santiago",
        "region": "Region Metropolitana",
        "country": "Chile",
        "lat": -33.45,
        "lon": -70.67,
        "tz_id": "America/Santiago",
        "localtime_epoch": 1721883278,
        "localtime": "2024-07-25 0:54"
    },
    "current": {
        "last_updated_epoch": 1721882700,
        "last_updated": "2024-07-25 00:45",
        "temp_c": 7.7,
        "temp_f": 45.9,
        "is_day": 0,
        "condition": {
            "text": "Overcast",
            "icon": "//cdn.weatherapi.com/weather/64x64/night/122.png",
            "code": 1009
        },
        "wind_mph": 2.2,
        "wind_kph": 3.6,
        "wind_degree": 95,
        "wind_dir": "E",
        "pressure_mb": 1024.0,
        "pressure_in": 30.23,
        "precip_mm": 0.0,
        "precip_in": 0.0,
        "humidity": 36,
        "cloud": 99,
        "feelslike_c": 7.9,
        "feelslike_f": 46.2,
        "windchill_c": 7.9,
        "windchill_f": 46.2,
        "heatindex_c": 7.7,
        "heatindex_f": 45.9,
        "dewpoint_c": -6.4,
        "dewpoint_f": 20.4,
        "vis_km": 10.0,
        "vis_miles": 6.0,
        "uv": 1.0,
        "gust_mph": 4.2,
        "gust_kph": 6.8
    }
}

We will be using from location: name, region, country.
And from current: temp_c, condition.Text and condition.Icon, wind_kph and feelslike_c

The steps to build this endpoint are:
	- create the data structs
	- make the function that returns the answer through a repository to be used by the usecases
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type apiRepository struct {
	apiKey string
}

func NewApiRepository(apiKey string) WeatherRepository {
	return &apiRepository{apiKey: apiKey}
}

// GetWeatherByCity fetches weather data sequentially
func (repo *apiRepository) GetWeatherByCity(city string) (WeatherData, error) {
	log.Println("Fetching weather data for city: ", city)
	url := fmt.Sprintf("https://weatherapi-com.p.rapidapi.com/current.json?q=%s", city)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WeatherData{}, err
	}

	// Set the necessary headers as per your cURL command
	req.Header.Set("x-rapidapi-host", "weatherapi-com.p.rapidapi.com")
	req.Header.Set("x-rapidapi-key", repo.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return WeatherData{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return WeatherData{}, err
	}

	var response struct {
		Location struct {
			Name    string `json:"name"`
			Region  string `json:"region"`
			Country string `json:"country"`
		} `json:"location"`
		Current struct {
			TempC     float64 `json:"temp_c"`
			Condition struct {
				Text string `json:"text"`
				Icon string `json:"icon"`
			} `json:"condition"`
			WindKPH    float64 `json:"wind_kph"`
			FeelsLikeC float64 `json:"feelslike_c"`
		} `json:"current"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return WeatherData{}, err
	}

	return WeatherData{
		City:       response.Location.Name,
		Region:     response.Location.Region,
		Country:    response.Location.Country,
		TempC:      response.Current.TempC,
		Condition:  response.Current.Condition.Text,
		IconURL:    fmt.Sprintf("https:%s", response.Current.Condition.Icon),
		WindKPH:    response.Current.WindKPH,
		FeelsLikeC: response.Current.FeelsLikeC,
	}, nil
}

// Placeholder for other methods - implement concurrency patterns here
func (repo *apiRepository) GetWeatherByCityWG(city string) (WeatherData, error) {
	// Implement using sync.WaitGroup
	return WeatherData{}, nil
}

func (repo *apiRepository) GetWeatherByCityChan(city string) (WeatherData, error) {
	// Implement using channels
	return WeatherData{}, nil
}

func (repo *apiRepository) GetWeatherByCityMx(city string) (WeatherData, error) {
	// Implement using mutexes
	return WeatherData{}, nil
}
