package repository

// NewsRepository is the interface for fetching news data
type NewsRepository interface {
	GetNewsByCity(city string) ([]NewsData, error)
	GetNewsByCityWG(city string) ([]NewsData, error)
	GetNewsByCityChan(city string) ([]NewsData, error)
	GetNewsByCityMx(city string) ([]NewsData, error)
}

// WeatherRepository is the interface that wraps the GetWeatherByCity method
type WeatherRepository interface {
	GetWeatherByCity(city string) (WeatherData, error)
	GetWeatherByCityWG(city string) (WeatherData, error)
	GetWeatherByCityChan(city string) (WeatherData, error)
	GetWeatherByCityMx(city string) (WeatherData, error)
}

// WeatherAPIResponse represents the full structure of the API response
type WeatherAPIResponse struct {
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

// WeatherData contains the simplified structure that will be used by the usecase layer
type WeatherData struct {
	City       string
	Region     string
	Country    string
	TempC      float64
	Condition  string
	IconURL    string
	WindKPH    float64
	FeelsLikeC float64
}

// NewsData represents the structured format of news articles
type NewsData struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Author      string   `json:"author"`
	Image       string   `json:"image"`
	Language    string   `json:"language"`
	Category    []string `json:"category"`
	Published   string   `json:"published"`
}
