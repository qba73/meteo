package main

import (
	"fmt"
	"log"

	"github.com/qba73/meteo"
)

func main() {
	// Get weather status for the given lat, lon:
	weather, err := meteo.GetWeather(53.2, -6.2)
	if err != nil {
		log.Println(err)
	}

	// Print out weather string.
	// Example: Lightrain 8.3°C
	fmt.Println(weather)
}
