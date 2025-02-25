// Example of running the program below
//
// ➜  meteo git:(master) ✗ go build -o example ./examples/concurrent/main.go
// ➜  meteo git:(master) ✗ ./example
// 0.61s Location: Lodz,PL, Weather: Cloudy 2.3°C
// 0.61s Location: Killarney,IE, Weather: Cloudy 12.0°C
// 0.61s Location: Vilnius,LT, Weather: Cloudy 4.3°C
// 0.61s Location: Berlin,DE, Weather: Partlycloudy_night 3.8°C
// 0.62s Location: Castlebar,IE, Weather: Cloudy 12.7°C
// 0.62s Location: Belfast,UK, Weather: Lightrain 10.5°C
// 0.62s Location: Dublin,IE, Weather: Cloudy 12.2°C
// 0.63s Location: Warsaw,PL, Weather: Partlycloudy_night 3.6°C
// 0.63s Location: London,UK, Weather: Cloudy 11.5°C
// 2.24s Location: Vienna,AT, Weather: Cloudy 4.6°C
// 2.24s elapsed
package main

import (
	"fmt"
	"time"

	"github.com/qba73/meteo"
)

func main() {
	start := time.Now()
	ch := make(chan string)

	locations := []string{
		"Vilnius,LT", "Dublin,IE", "London,UK", "Berlin,DE",
		"Belfast,UK", "Castlebar,IE", "Killarney,IE",
		"Warsaw,PL", "Lodz,PL", "Vienna,AT",
	}

	for _, loc := range locations {
		go getWeather(loc, ch)
	}

	for range locations {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func getWeather(location string, ch chan<- string) {
	start := time.Now()

	weather, err := meteo.GetWeather(location)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	sec := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs Location: %s, Weather: %s", sec, location, weather)
}
