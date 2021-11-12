package meteo

import (
	"fmt"
	"os"
	"strings"

	"github.com/qba73/meteo/geonames"
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

// RunCLI is a main function that runs the cli machinery.
func RunCLI() {
	uname := os.Getenv("GEO_USERNAME")
	resolver, err := geonames.NewClient(uname)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	c, err := NewYrClient(resolver)
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
