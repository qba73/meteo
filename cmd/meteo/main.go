package main

import (
	"os"

	"github.com/qba73/meteo"
)

func main() {
	os.Exit(meteo.RunWeatherCLI())
}
