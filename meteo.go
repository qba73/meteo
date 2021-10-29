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

func RunCLI() {
	c := NewNorwayClient()
	w, err := c.GetForecast(53.2, -6.2)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(w)
}
