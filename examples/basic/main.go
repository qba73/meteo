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

	// Print weather.
	fmt.Println(weather)

	// Print summary and temperature in Celcius.
	fmt.Println(weather.Summary)
	fmt.Println(weather.Temp)
}
