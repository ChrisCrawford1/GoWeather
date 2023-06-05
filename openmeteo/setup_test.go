package openmeteo

import (
	"net/http"
	"os"
	"testing"
)

var geocodingJsonReturn = `
{
	"results": [{
		"id": 6167865,
		"name": "Toronto",
		"latitude": 43.70011,
		"longitude": -79.4163,
		"elevation": 175,
		"feature_code": "PPLA",
		"country_code": "CA",
		"admin1_id": 6093943,
		"timezone": "America/Toronto",
		"population": 2600000,
		"country_id": 6251999,
		"country": "Canada",
		"admin1": "Ontario"
	}],
	"generationtime_ms": 0.502944
}`

var geocodingEmptyResultsJsonReturn = `
{
	"results": [],
	"generationtime_ms": 0.502944
}`

var locationWeatherJsonReturn = `
{
	"latitude": 43.70455,
	"longitude": -79.4046,
	"generationtime_ms": 0.14293193817138672,
	"utc_offset_seconds": -14400,
	"timezone": "America/New_York",
	"timezone_abbreviation": "EDT",
	"elevation": 175,
	"current_weather": {
		"temperature": 25.4,
		"windspeed": 10.5,
		"winddirection": 354,
		"weathercode": 3,
		"is_day": 0,
		"time": "2023-06-02T21:00"
	}
}`

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type RoundTripFunc func(req *http.Request) *http.Response

func (rTrip RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return rTrip(req), nil
}

func NewTestHttpClient(rTrip RoundTripFunc) *http.Client {
	// For testing purposes we want to mock out the transport layer / call
	// so as to not ping any real world services when the tests run.
	return &http.Client{
		Transport: rTrip,
	}
}
