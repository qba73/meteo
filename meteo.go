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

type NameResolver interface {
	GetCoordinates(placeName, country string) (Place, error)
}

func RunCLI() {
	uname := os.Getenv("GEO_USERNAME")
	resolver, err := NewWikipediaClient(uname)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	c, err := NewNorwayClient(resolver)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	w, err := c.GetForecast("Castlebar", "IE")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(w)
}
