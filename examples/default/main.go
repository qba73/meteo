package main

import (
	"fmt"
	"log"

	"github.com/qba73/meteo"
)

func main() {
	// Get weather status for city Castlebar in Ireland:
	weather, err := meteo.GetWeather("Castlebar", "IE")
	if err != nil {
		log.Println(err)
	}

	// Print out weather string.
	// Example: Lightrain 8.3°C
	fmt.Println(weather)
}
