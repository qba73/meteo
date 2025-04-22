[![Go Reference](https://pkg.go.dev/badge/github.com/qba73/meteo.svg)](https://pkg.go.dev/github.com/qba73/meteo)
![GitHub](https://img.shields.io/github/license/qba73/meteo)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/qba73/meteo)


# meteo

`meteo` is a Go client library for the weather and meteorological forecast from [Yr](https://www.yr.no/en).

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
	// Export GEO_USERNAME Env Var (you registered at geonames.org)

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

# Installation
```
$ go install github.com/qba73/meteo/cmd/meteo@latest
```
