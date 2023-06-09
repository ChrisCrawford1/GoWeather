# Go Weather

[![Go](https://github.com/ChrisCrawford1/GoWeather/actions/workflows/go.yml/badge.svg)](https://github.com/ChrisCrawford1/GoWeather/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/ChrisCrawford1/GoWeather/branch/main/graph/badge.svg?token=2LGbW1rftR)](https://codecov.io/gh/ChrisCrawford1/GoWeather)


## About
A small CLI application that can be used to fetch weather for any city. This was mainly built as a tool for myself that I use frequently, especially in regard to the air quality for my alergies during the summer months.

_UV index and air quality index will be added soon, as will pollen levels_

_(Only available for Europe during pollen season however.)_
---


**Current features**
* Temperature (Celcius)
* Wind direction
* Wind speed


**City name entry view**

![initial](https://github.com/ChrisCrawford1/GoWeather/assets/44769623/ffaafae0-3450-4b7a-ae7e-da00129ea6f2)

**Weather info view**

![info](https://github.com/ChrisCrawford1/GoWeather/assets/44769623/7c1b4e02-5f1b-49d7-b518-9e44258f41d2)


**Planned features**
* UV Index
* US Air Quality Index
* UI theme updates / better colouring


## Built with
* Golang - 1.20
* Bubbletea - 0.24
* OpenMeteo

---


## Local Development
To run the app locally, clone the repository into your desired location. 
```bash
    cd $DIR
    make dev
```

To run tests with no coverage report
```bash
    make test
```

To run with generated coverage report and display in browser
```bash
    make generate-coverage && make show-coverage
```
Of course these commands can be run one at a time if you wish.
