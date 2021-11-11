package meteo

import (
	"fmt"
	"os"
	"strings"
)

const (
	userAgent = "Meteo/0.1 https://github.com/qba73/meteo"
)

// Weather represents weather conditions
// in a geographical region.
type Weather struct {
	Summary string
	Temp    float64
}

// String implements stringer interface.
func (w Weather) String() string {
	return fmt.Sprintf("%s %.1f°C", strings.Title(w.Summary), w.Temp)
}

// NameResolver interface is used by an Meteo Client
// to obtain geo coordinates for given place located in
// a country identified by country id.
type NameResolver interface {
	// GetCoordinates takes place and country code
	// and returns geo information like lat and lng.
	GetCoordinates(placeName, country string) (Place, error)
}

// RunCLI is a main function that runs the cli machinery.
func RunCLI() {
	uname := os.Getenv("GEO_USERNAME")
	c, err := NewYrWeatherClient(
		WithWikipediaGeoResolver(),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	c, err := NewNorwayClient(resolver)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	forecast, err := c.GetForecast("Castlebar", "IE")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(w)
}
