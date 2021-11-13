package meteo

import (
	"fmt"
	"os"
	"strings"
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
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s LOCATION\n\nExmple: %[1]s London,UK\n", os.Args[0])
		os.Exit(1)
	}
	location := strings.Join(os.Args[1:], " ")
	w, err := GetWeather(location)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(w)
}
