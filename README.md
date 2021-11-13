![Go](https://github.com/qba73/meteo/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/qba73/meteo)](https://goreportcard.com/report/github.com/qba73/meteo)
[![Maintainability](https://api.codeclimate.com/v1/badges/4afc34a390da95ed9327/maintainability)](https://codeclimate.com/github/qba73/meteo/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/4afc34a390da95ed9327/test_coverage)](https://codeclimate.com/github/qba73/meteo/test_coverage)
![GitHub](https://img.shields.io/github/license/qba73/meteo)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/qba73/meteo)


# meteo

Meteo is a Go client library for the weather and meteorological forecast from [Yr](https://www.yr.no/en).

> Weather forecast from Yr, delivered by the Norwegian Meteorological Institute and NRK.

# Usage

You must register your user agent string in the [YR.NO service](https://developer.yr.no/doc/TermsOfService/) and your user name at the [GeoNames.org](https://www.geonames.org/login).


```go
package main

import (
	"fmt"

	"github.com/qba73/meteo"
)

func main() {
	// Export GEO_USERNAME Env Var - username you registered at geonames.org

	// Get current weather for given location.
	weather, err := meteo.GetWeather("Vilnius,LT")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(weather)
	// Cloudy 4.2°C

	fmt.Println(weather.Summary)
	// cloudy

	fmt.Println(weather.Temp)
	// 4.2
}
```

The code sample below shows a basic example of how the meteo package can fetch weather statuses concurrently.

```go
package main

import (
	"fmt"
	"time"

	"github.com/qba73/meteo"
)

func main() {
	start := time.Now()
	ch := make(chan string)

	locations := []string{
		"Vilnius,LT", "Dublin,IE", "London,UK", "Berlin,DE",
		"Belfast,UK", "Castlebar,IE", "Killarney,IE",
		"Warsaw,PL", "Lodz,PL", "Vienna,AT"}

	for _, loc := range locations {
		go getWeather(loc, ch)
	}

	for range locations {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func getWeather(location string, ch chan<- string) {
	start := time.Now()

	weather, err := meteo.GetWeather(location)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	sec := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs Location: %s, Weather: %s", sec, location, weather)
}
```
Build the binary:
```bash
➜  meteo git:(master) ✗ go build -o example ./examples/concurrent/main.go
```
Run the binary:
```
➜  meteo git:(master) ✗ ./example
0.61s Location: Lodz,PL, Weather: Cloudy 2.3°C
0.61s Location: Killarney,IE, Weather: Cloudy 12.0°C
0.61s Location: Vilnius,LT, Weather: Cloudy 4.3°C
0.61s Location: Berlin,DE, Weather: Partlycloudy_night 3.8°C
0.62s Location: Castlebar,IE, Weather: Cloudy 12.7°C
0.62s Location: Belfast,UK, Weather: Lightrain 10.5°C
0.62s Location: Dublin,IE, Weather: Cloudy 12.2°C
0.63s Location: Warsaw,PL, Weather: Partlycloudy_night 3.6°C
0.63s Location: London,UK, Weather: Cloudy 11.5°C
2.24s Location: Vienna,AT, Weather: Cloudy 4.6°C
2.24s elapsed
```

# Installation
```
$ go get github.com/qba73/meteo
```
