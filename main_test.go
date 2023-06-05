package main

import (
	"errors"
	"fmt"
	"goweather/openmeteo"
	"reflect"
	"testing"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TestViewCursorType(t *testing.T) {
	textInputField := textinput.NewModel()
	var testModel = Model{
		typing:    true,
		textInput: textInputField,
	}

	returnedCmd := testModel.Init()
	returnedType := reflect.TypeOf(returnedCmd).String()

	if returnedType != "tea.Cmd" {
		t.Errorf("Unexpected init return type, was expecting tea.Cmd, got %s", returnedType)
	}
}

func TestViewTypingStringReturn(t *testing.T) {
	textInputField := textinput.NewModel()
	var testModel = Model{
		typing:    true,
		textInput: textInputField,
	}

	returnedString := testModel.View()

	expectedString := fmt.Sprintf(
		"Enter City Name: \n%s\n%s",
		testModel.textInput.View(),
		"(Ctrl+C to quit | esc to clear input)",
	)

	if returnedString != expectedString {
		t.Errorf("Unexpected string returned from view, got: %s", returnedString)
	}
}

func TestViewLoadingStringReturn(t *testing.T) {
	loadingSpinner := spinner.NewModel()
	loadingSpinner.Spinner = spinner.Points

	textInputField := textinput.NewModel()
	textInputField.SetValue("Toronto")

	var testModel = Model{
		textInput:      textInputField,
		loading:        true,
		loadingSpinner: loadingSpinner,
	}

	returnedString := testModel.View()
	// NOTE: in real use the spinner is animated, in tests it appears as static however.
	expectedString := fmt.Sprintf(
		"Fetching weather for %s please wait... %s", testModel.textInput.Value(), testModel.loadingSpinner.View(),
	)

	if returnedString != expectedString {
		t.Errorf("Got %s", returnedString)
	}
}

func TestViewErrorStringReturn(t *testing.T) {
	returnedError := errors.New("Unexpected status code received from OpenMeteo Status=500")
	var testModel = Model{
		err: returnedError,
	}

	returnedString := testModel.View()
	expectedString := fmt.Sprintf("%s\n", returnedError.Error())

	if returnedString != expectedString {
		t.Errorf("Unexpected error returned %s", returnedError.Error())
	}
}

func TestViewWeatherInfoStringReturn(t *testing.T) {
	var testModel = Model{
		city:    "Toronto",
		country: "Canada",
		weather: openmeteo.CurrentWeather{
			Latitude:             52.42342,
			Longitude:            -47.23453,
			GenerationtimeMs:     10.0,
			UtcOffsetSeconds:     1,
			Timezone:             "America/Toronto",
			TimezoneAbbreviation: "EST",
			Elevation:            145.0,
			WeatherInfo: openmeteo.WeatherInfoVariables{
				Temperature:   30.2,
				Windspeed:     12.1,
				Winddirection: 112.0,
				Weathercode:   1.1,
				IsDay:         0,
				Time:          "2023-06-02T20:00",
			},
		},
	}

	currentWeather := testModel.weather.WeatherInfo
	returnedString := testModel.View()
	expectedString := fmt.Sprintf(
		"Weather Information for %s, %s @ %s\nCurrent Temp %.1fc\nWind %.1fº @ %.1fkts\n%s",
		testModel.city,
		testModel.country,
		currentWeather.GetFormattedTime(),
		currentWeather.Temperature,
		currentWeather.Winddirection,
		currentWeather.Windspeed,
		"(Ctrl+C to quit | esc to clear input)",
	)

	if returnedString != expectedString {
		t.Error("Unexpected weather output received")
	}
}

func TestUpdateExitCase(t *testing.T) {
	var testModel = Model{
		typing: true,
	}

	testMsg := tea.KeyMsg{Type: tea.KeyCtrlC}

	_, cmd := testModel.Update(testMsg)

	switch msg := cmd().(type) {
	case tea.QuitMsg:
		// Do nothing, as this is the pass case
	default:
		// Variable is not of type QuitMsg
		// Handle other types here
		t.Error("Expected tea.QuitMsg from this action, received %w", msg)
	}
}

func TestUpdateEscCase(t *testing.T) {
	textInputField := textinput.NewModel()
	textInputField.Focus()
	textInputField.SetValue("Toronto")

	var testModel = Model{
		typing:    true,
		textInput: textInputField,
	}

	testMsg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := testModel.Update(testMsg)

	if cmd != nil {
		t.Error("Expected a nil value for command, non nil value received")
	}

	// TODO: Update test to ensure the input field is set to empty and renders empty also

	// expected := fmt.Sprintf(
	// 	"Enter City Name: \n>\n%s",
	// 	"(Ctrl+C to quit | esc to clear input)",
	// )
	// if model.View() != expected {
	// 	t.Errorf("Expected %s, got %s", expected, model.View())
	// }
}

func TestUpdateFetchedWeatherErrorCase(t *testing.T) {
	fetchedWeather := FetchedWeather{
		Err: errors.New("Something went wrong fetching the data"),
	}
	var testModel = Model{
		loading: true,
	}

	_, cmd := testModel.Update(fetchedWeather)

	if cmd != nil {
		t.Error("Expected a nil value for command, non nil value received")
	}
}
