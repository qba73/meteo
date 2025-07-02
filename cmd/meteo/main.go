package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/qba73/meteo"
)

func main() {
	os.Exit(run(os.Stdout, os.Stderr))
}

func run(w, ew io.Writer) int {
	if len(os.Args) < 2 {
		fmt.Fprintf(ew, "Usage: %s LOCATION\n\nExmple: %[1]s London,UK\n", os.Args[0])
		return 1
	}
	location := strings.Join(os.Args[1:], " ")
	weather, err := meteo.GetWeather(location)
	if err != nil {
		fmt.Fprintln(ew, err)
		return 1
	}
	fmt.Fprintln(w, weather)
	return 0
}
