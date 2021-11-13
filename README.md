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

# Installation
```
$ go get github.com/qba73/meteo
```
