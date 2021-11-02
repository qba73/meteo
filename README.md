![Go](https://github.com/qba73/meteo/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/qba73/meteo)](https://goreportcard.com/report/github.com/qba73/meteo)
[![Maintainability](https://api.codeclimate.com/v1/badges/4afc34a390da95ed9327/maintainability)](https://codeclimate.com/github/qba73/meteo/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/4afc34a390da95ed9327/test_coverage)](https://codeclimate.com/github/qba73/meteo/test_coverage)


# meteo

Meteo is a Go client library for the weather and meteorological forecast from [Yr](https://www.yr.no/en).

Disclaimer:

Weather forecast from Yr, delivered by the Norwegian Meteorological Institute and NRK.

# Usage

## Preconditions

You must register your user agent string in the [YR.NO service](https://developer.yr.no/doc/TermsOfService/) and your user name in the [GeoNames service](https://www.geonames.org/login) to use the package.


## Installation
```
$ go get git@github.com:qba73/meteo.git
```

## Default

Export ```GEO_USERNAME``` env var that you registered with GeoNames.org.

Example:
```
$ export GEO_USERNAME=Jane123
```
Use the ```meteo``` package in your application:
```go
$ package main

import (
	"fmt"
	"log"
	"github.com/qba73/meteo"
)

func main() {
	// Get weather status for Castlebar in Ireland:
	weather, err := meteo.GetWeather("Castlebar", "IE")
	if err != nil {
		log.Println(err)
	}
	// Print out weather string.
	// Example: Lightrain 8.3°C
	fmt.Println(weather)
}
```

