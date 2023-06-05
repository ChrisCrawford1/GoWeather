# Go Weather

<!-- [![Go](https://github.com/ChrisCrawford1/Command/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/ChrisCrawford1/Command/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/ChrisCrawford1/JobTrack/branch/main/graph/badge.svg?token=YPj3MteVDP)]() -->


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

![City name entry view](./images/initial.png)

**Weather info view**

![City name entry view](./images/info.png)


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
