package main

import (
	"os"
)

func main() {
	// Get coordinates using a default Geo client
	user := os.Getenv("GEO_USERNAME")
	_ = user
	/*
		coord, err := geonames.
		if err != nil {
			println(err)
		}

		fmt.Printf("Lat: %.2f, Lng: %.2f for %s in country %s\n", coord.Lat, coord.Lng, coord.PlaceName, coord.CountryCode)
		// It returns:
		// Lat: 53.85, Lng: -9.30 for Castlebar in country IE
	*/
}
