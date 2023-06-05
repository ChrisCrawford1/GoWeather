package openmeteo

import (
	"bytes"
	"io"
	"net/http"
	"testing"
)

var weatherInfoVars = WeatherInfoVariables{
	Temperature:   21.0,
	Windspeed:     10.1,
	Winddirection: 114.0,
	Weathercode:   2.2,
	IsDay:         0,
	Time:          "2023-06-02T20:00",
}

func TestClientGetFormattedTimeWithValidTimeString(t *testing.T) {
	formattedTime := weatherInfoVars.GetFormattedTime()

	if formattedTime != "Jun 2, 2023 20:00" {
		t.Error("Unexpected time format found")
	}
}

func TestClientGetFormattedTimeWithInvalidTimeString(t *testing.T) {
	weatherInfoVars.Time = "I am broken"

	formattedTime := weatherInfoVars.GetFormattedTime()

	if formattedTime != "Time error" {
		t.Errorf("Expected a string of Time error, got instead %s", formattedTime)
	}
}

func TestClientGetGeocodingLocationData(t *testing.T) {
	testClient := NewTestHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(geocodingJsonReturn)),
			Header:     make(http.Header),
		}
	})

	client := Client{
		HTTPClient: testClient,
	}

	result, err := client.GetGeocodingForLocation("Toronto")

	if err != nil {
		t.Errorf("Unexpected error received %s", err.Error())
	}

	if len(result.Results) == 0 {
		t.Errorf("Expected geocoding results to be populated, got length of %d", len(result.Results))
	}
}

func TestClientGetGeocodingLocationDataNonOkResult(t *testing.T) {
	testClient := NewTestHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBufferString("")),
			Header:     make(http.Header),
		}
	})

	client := Client{
		HTTPClient: testClient,
	}

	_, err := client.GetGeocodingForLocation("Toronto")

	if err == nil {
		t.Error("Was expecting status code error, none received")
	}
}

func TestClientGetGeocodingLocationDataNoResults(t *testing.T) {
	testClient := NewTestHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(geocodingEmptyResultsJsonReturn)),
			Header:     make(http.Header),
		}
	})

	client := Client{
		HTTPClient: testClient,
	}

	_, err := client.GetGeocodingForLocation("vlahbah")
	expectedError := "Couldn't find a match for vlahbah, please try again"

	if err.Error() != expectedError {
		t.Errorf("Expected error for no results, got %s", err)
	}
}

func TestClientGetGeocodingLocationRoundTripError(t *testing.T) {
	testClient := NewTestHttpClientWithError(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("")),
			Header:     make(http.Header),
		}
	})

	client := Client{
		HTTPClient: testClient,
	}

	_, err := client.GetGeocodingForLocation("vlahbah")
	expectedError := `Error fetching geocoding data from OpenMeteo: Get "https://geocoding-api.open-meteo.com/v1/search?name=vlahbah&count=1": Something went wrong`

	if err.Error() != expectedError {
		t.Errorf("Expected error for fetching data, got %s", err)
	}
}

func TestClientGetWeatherInfoFromCoords(t *testing.T) {
	testClient := NewTestHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(locationWeatherJsonReturn)),
			Header:     make(http.Header),
		}
	})

	client := Client{
		HTTPClient: testClient,
	}

	_, err := client.GetWeatherInfoFromCoords(43.70455, -79.4046)

	if err != nil {
		t.Errorf("Unexpected error received %s", err.Error())
	}
}

func TestClientGetWeatherInfoFromCoordsNonOkResult(t *testing.T) {
	testClient := NewTestHttpClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBufferString("")),
			Header:     make(http.Header),
		}
	})

	client := Client{
		HTTPClient: testClient,
	}

	_, err := client.GetWeatherInfoFromCoords(43.70455, -79.4046)

	if err == nil {
		t.Error("Was expecting status code error, none received")
	}
}

func TestClientGetWeatherInfoFromCoordsRoundTripError(t *testing.T) {
	testClient := NewTestHttpClientWithError(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("")),
			Header:     make(http.Header),
		}
	})

	client := Client{
		HTTPClient: testClient,
	}

	_, err := client.GetWeatherInfoFromCoords(43.70455, -79.4046)
	expectedError := `Error fetching weather data from OpenMeteo: Get "https://api.open-meteo.com/v1/forecast?latitude=43.7045&longitude=-79.4046&current_weather=true&timezone=EST": Something went wrong`

	if err.Error() != expectedError {
		t.Errorf("Expected error for fetching data, got %s", err)
	}
}
