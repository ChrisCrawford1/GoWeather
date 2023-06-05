package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"GoWeather/openmeteo"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type FetchedWeather struct {
	Weather openmeteo.CurrentWeather
	Err     error
	City    string
	Country string
}

type Model struct {
	openMeteoClient *openmeteo.Client
	weather         openmeteo.CurrentWeather
	city            string
	country         string
	textInput       textinput.Model
	loadingSpinner  spinner.Model
	typing          bool
	loading         bool
	err             error
}

func main() {
	textInputField := textinput.NewModel()
	textInputField.Focus()

	loadingSpinner := spinner.NewModel()
	loadingSpinner.Spinner = spinner.Points

	initModel := &Model{
		openMeteoClient: &openmeteo.Client{HTTPClient: http.DefaultClient},
		textInput:       textInputField,
		loadingSpinner:  loadingSpinner,
		typing:          true,
	}

	err := tea.NewProgram(initModel, tea.WithAltScreen()).Start()

	if err != nil {
		log.Panic(err)
	}
}

func (m Model) getWeatherFromUserInput(cityName string) tea.Cmd {
	return func() tea.Msg {
		geocoding, err := m.openMeteoClient.GetGeocodingForLocation(cityName)

		if err != nil {
			return FetchedWeather{Err: err}
		}

		latitude := geocoding.Results[0].Latitude
		longitude := geocoding.Results[0].Longitude
		locationCurrentWeather, err := m.openMeteoClient.GetWeatherInfoFromCoords(latitude, longitude)

		if err != nil {
			return FetchedWeather{Err: err}
		}

		return FetchedWeather{
			Weather: locationCurrentWeather,
			City:    geocoding.Results[0].Name,
			Country: geocoding.Results[0].Country,
		}
	}
}

func (m Model) Init() tea.Cmd {
	// Not really needed, just makes things look a little more
	// interesting as opposed to the standard empty cursor.
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Let the user reuse the app if they want without having to ctrl+c and rerun.
			m.typing = true
			m.err = nil

			// Delete the previous input so they can start again from scratch.
			m.textInput.Reset()

			return m, nil
		case "enter":
			cityName := strings.TrimSpace(m.textInput.Value())

			if cityName != "" {
				m.typing = false
				m.loading = true

				return m, tea.Batch(
					spinner.Tick,
					m.getWeatherFromUserInput(cityName),
				)
			}
		}
	case FetchedWeather:
		m.loading = false
		if err := msg.Err; err != nil {
			m.err = err
			return m, nil
		}

		m.weather = msg.Weather
		m.city = msg.City
		m.country = msg.Country
		return m, nil
	}

	if m.typing {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)

		return m, cmd
	}

	if m.loading {
		var cmd tea.Cmd
		m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.typing {
		return fmt.Sprintf(
			"Enter City Name: \n%s\n%s",
			m.textInput.View(),
			"(Ctrl+C to quit | esc to clear input)",
		)
	}

	if m.loading {
		return fmt.Sprintf("Fetching weather for %s please wait... %s", m.textInput.Value(), m.loadingSpinner.View())
	}

	if err := m.err; err != nil {
		return fmt.Sprintf("%s\n", err)
	}

	currentWeather := m.weather.WeatherInfo

	return fmt.Sprintf(
		"Weather Information for %s, %s @ %s\nCurrent Temp %.1fc\nWind %.1fº @ %.1fkts\n%s",
		m.city, m.country, currentWeather.GetFormattedTime(), currentWeather.Temperature, currentWeather.Winddirection, currentWeather.Windspeed,
		"(Ctrl+C to quit | esc to clear input)",
	)
}
