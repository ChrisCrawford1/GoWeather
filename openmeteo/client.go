package openmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	HTTPClient *http.Client
}

type Geocoding struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Elevation   int     `json:"-"`
	FeatureCode string  `json:"-"`
	CountryCode string  `json:"country_code"`
	Admin1ID    int     `json:"-"`
	Timezone    string  `json:"timezone"`
	Population  int     `json:"population"`
	CountryID   int     `json:"-"`
	Country     string  `json:"country"`
	Admin1      string  `json:"-"`
}

type GeocodingResults struct {
	Results        []Geocoding `json:"results"`
	GenerationTime float64     `json:"generationtime_ms"`
}

type WeatherInfoVariables struct {
	Temperature   float64 `json:"temperature"`
	Windspeed     float64 `json:"windspeed"`
	Winddirection float32 `json:"winddirection"`
	Weathercode   float32 `json:"weathercode"`
	IsDay         int     `json:"is_day"`
	Time          string  `json:"time"`
}

type CurrentWeather struct {
	Latitude             float64              `json:"latitude"`
	Longitude            float64              `json:"longitude"`
	GenerationtimeMs     float64              `json:"generationtime_ms"`
	UtcOffsetSeconds     int                  `json:"utc_offset_seconds"`
	Timezone             string               `json:"timezone"`
	TimezoneAbbreviation string               `json:"timezone_abbreviation"`
	Elevation            float32              `json:"elevation"`
	WeatherInfo          WeatherInfoVariables `json:"current_weather"`
}

func (weatherVariables *WeatherInfoVariables) GetFormattedTime() string {
	t, err := time.Parse("2006-01-02T15:04", weatherVariables.Time)

	var formattedTime string

	if err != nil {
		formattedTime = "Time error"
	} else {
		formattedTime = t.Format("Jan 2, 2006 15:04")
	}

	return formattedTime
}

func (client *Client) GetGeocodingForLocation(cityName string) (GeocodingResults, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1", cityName),
		nil,
	)

	if err != nil {
		return GeocodingResults{}, fmt.Errorf("Could not create a fetch city geocode request: %w", err)
	}

	resp, err := client.HTTPClient.Do(req)

	if err != nil {
		return GeocodingResults{}, fmt.Errorf("Error fetching geocoding data from OpenMeteo: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GeocodingResults{}, fmt.Errorf("Unexpected status code received from OpenMeteo Status=%d", resp.StatusCode)
	}

	var geocodingInfo GeocodingResults
	err = json.NewDecoder(resp.Body).Decode(&geocodingInfo)

	if err != nil {
		return GeocodingResults{}, fmt.Errorf("Error while decoding Geocoding JSON response: %w", err)
	}

	if len(geocodingInfo.Results) == 0 {
		// If we cant find a geolocation then return an error and let the user try again
		return GeocodingResults{}, fmt.Errorf("Couldn't find a match for %s, please try again", cityName)
	}

	return geocodingInfo, nil
}

func (client *Client) GetWeatherInfoFromCoords(latitude float64, longitude float64) (CurrentWeather, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current_weather=true&timezone=EST", latitude, longitude),
		nil,
	)

	if err != nil {
		return CurrentWeather{}, fmt.Errorf("Could not create a fetch weather request: %w", err)
	}

	resp, err := client.HTTPClient.Do(req)

	if err != nil {
		return CurrentWeather{}, fmt.Errorf("Error fetching weather data from OpenMeteo: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return CurrentWeather{}, fmt.Errorf("Unexpected status code received from OpenMeteo Status=%d", resp.StatusCode)
	}

	var currentWeather CurrentWeather
	err = json.NewDecoder(resp.Body).Decode(&currentWeather)

	if err != nil {
		return CurrentWeather{}, fmt.Errorf("Error while decoding Geocoding JSON response: %w", err)
	}

	return currentWeather, nil
}
