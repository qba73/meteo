package meteo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/qba73/geonames"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	libVersion = "0.0.1"
)

type forecastResponseCompact struct {
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Meta struct {
			UpdatedAt time.Time `json:"updated_at"`
			Units     struct {
				AirPressureAtSeaLevel string `json:"air_pressure_at_sea_level"`
				AirTemperature        string `json:"air_temperature"`
				CloudAreaFraction     string `json:"cloud_area_fraction"`
				PrecipitationAmount   string `json:"precipitation_amount"`
				RelativeHumidity      string `json:"relative_humidity"`
				WindFromDirection     string `json:"wind_from_direction"`
				WindSpeed             string `json:"wind_speed"`
			} `json:"units"`
		} `json:"meta"`
		Timeseries []struct {
			Time time.Time `json:"time"`
			Data struct {
				Instant struct {
					Details struct {
						AirPressureAtSeaLevel float64 `json:"air_pressure_at_sea_level"`
						AirTemperature        float64 `json:"air_temperature"`
						CloudAreaFraction     float64 `json:"cloud_area_fraction"`
						RelativeHumidity      float64 `json:"relative_humidity"`
						WindFromDirection     float64 `json:"wind_from_direction"`
						WindSpeed             float64 `json:"wind_speed"`
					} `json:"details"`
				} `json:"instant"`
				Next12Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
				} `json:"next_12_hours"`
				Next1Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						PrecipitationAmount float64 `json:"precipitation_amount"`
					} `json:"details"`
				} `json:"next_1_hours"`
				Next6Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						PrecipitationAmount float64 `json:"precipitation_amount"`
					} `json:"details"`
				} `json:"next_6_hours"`
			} `json:"data"`
		} `json:"timeseries"`
	} `json:"properties"`
}

type CurrentWeather struct {
	UpdatedAt          time.Time
	Time               time.Time
	PressureAtSeaLevel float64
	Temperature        float64
	Precipitation      float64
}

type Forecast struct {
	UpdatedAt time.Time
	Hourly    []HourlyForecast
}

type HourlyForecast struct {
	Time                time.Time
	AirPressure         float64 // in hPa
	AirTemperature      float64 // in Celsius
	CloudAreaFraction   float64 // in "%""
	RelativeHumidity    float64 // in "%"
	WindFromDirection   float64 // in "degrees"
	WindSpeed           float64 // in "m/s"
	PrecipitationAmount float64 // in "mm"
	Summary             string
}

// Location represents geo coordinates.
type Location struct {
	Lat  float64
	Long float64
}

func resolve(ctx context.Context, location string) (Location, error) {
	geonamesUser := os.Getenv("GEONAMES_USER")
	if geonamesUser == "" {
		log.Fatal("Please set environmental variable GEONAMES_USER.")
	}
	resolver, err := geonames.NewClient(geonamesUser)
	if err != nil {
		return Location{}, err
	}
	place, country, err := toPlaceAndCountry(location)
	if err != nil {
		return Location{}, err
	}
	names, err := resolver.GetPlace(place, country, 1)
	if err != nil {
		return Location{}, err
	}
	if len(names) < 1 {
		return Location{}, fmt.Errorf("unable to resolve location: place %s, country %s", place, country)
	}
	return Location{
		Lat:  names[0].Position.Lat,
		Long: names[0].Position.Long,
	}, nil
}

func toPlaceAndCountry(location string) (string, string, error) {
	bits := strings.Split(location, ",")
	if len(bits) < 2 {
		return "", "", fmt.Errorf("parsing place and country from location %s", location)
	}
	return bits[0], bits[1], nil
}

var summary map[string]string = map[string]string{
	"rain":                   "rain",
	"heavyrain":              "heavy rain",
	"lightrain":              "light rain",
	"cloudy":                 "cloudy",
	"heavyrainshowers_night": "heavy rain showers",
	"heavyrainshowers_day":   "heavy rain showers",
	"rainshowers_day":        "rain showers",
	"rainshowers_night":      "rain showers",
	"lightrainshowers_day":   "light showers",
	"lightrainshowers_night": "light showers",
	"partlycloudy_day":       "partly cloudy",
	"partlycloudy_night":     "partly cloudy",
	"fair_day":               "fair",
	"fair_night":             "fair",
	"fog":                    "fog",
	"clearsky_night":         "clear sky",
	"clearsky_day":           "clear sky",
}

type option func(*Client) error

func WithUserAgent(ua string) option {
	return func(c *Client) error {
		if ua == "" {
			return errors.New("nil user agent")
		}
		c.UserAgent = ua
		return nil
	}
}

func WithBaseURL(u string) option {
	return func(c *Client) error {
		if u == "" {
			return errors.New("nil base URL")
		}
		c.BaseURL = u
		return nil
	}
}

func WithHTTPClient(hc *http.Client) option {
	return func(c *Client) error {
		if hc == nil {
			return errors.New("nil http client")
		}
		c.HTTPClient = hc
		return nil
	}
}

func WithResolver(rs func(context.Context, string) (Location, error)) option {
	return func(c *Client) error {
		c.Resolve = rs
		return nil
	}
}

// Client represents a weather client
// for the Norwegian Meteorological Institute.
type Client struct {
	UserAgent  string
	BaseURL    string
	HTTPClient *http.Client
	Resolve    func(context.Context, string) (Location, error)
}

// NewClient knows how to construct a new default client.
func NewClient(opts ...option) (*Client, error) {
	c := Client{
		UserAgent:  "Meteo/" + libVersion + "https://github.com/qba73/meteo",
		BaseURL:    "https://api.met.no",
		HTTPClient: http.DefaultClient,
		Resolve:    resolve,
	}
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}
	return &c, nil
}

// GetWeather returns current weather for given place.
//
// Place string should have format: "<place-name>,<country-code>",
// for example: "London,UK", "Dublin,IE", "Paris,FR", "Warsaw,PL".
func (c Client) GetWeather(ctx context.Context, place string) (Weather, error) {
	location, err := c.Resolve(ctx, place)
	if err != nil {
		return Weather{}, err
	}
	return c.weather(ctx, location)
}

// GetWeatherForCoordinates returns current weather for a place
// with given coordinates (lat, long)
func (c Client) GetWeatherForCoordinates(ctx context.Context, lat, long float64) (Weather, error) {
	l := Location{
		Lat:  lat,
		Long: long,
	}
	return c.weather(ctx, l)

}

func (c Client) weather(ctx context.Context, location Location) (Weather, error) {
	u := fmt.Sprintf("%s/weatherapi/locationforecast/2.0/compact?lat=%.2f&lon=%.2f", c.BaseURL, location.Lat, location.Long)

	var nf forecastResponseCompact
	if err := c.get(ctx, u, &nf); err != nil {
		return Weather{}, err
	}

	if len(nf.Properties.Timeseries) < 1 {
		return Weather{}, fmt.Errorf("invalid response %+v", nf)
	}

	w := Weather{
		Summary: nf.Properties.Timeseries[0].Data.Next1Hours.Summary.SymbolCode,
		Temp:    nf.Properties.Timeseries[0].Data.Instant.Details.AirTemperature,
	}
	return w, nil
}

func (c Client) forecast(ctx context.Context, location Location) (Forecast, error) {
	u := fmt.Sprintf("%s/weatherapi/locationforecast/2.0/compact?lat=%.2f&lon=%.2f", c.BaseURL, location.Lat, location.Long)

	var nf forecastResponseCompact
	if err := c.get(ctx, u, &nf); err != nil {
		return Forecast{}, err
	}

	if len(nf.Properties.Timeseries) < 1 {
		return Forecast{}, fmt.Errorf("invalid response %+v", nf)
	}

	w := Forecast{
		UpdatedAt: nf.Properties.Meta.UpdatedAt,
	}
	return w, nil
}

func (c Client) get(ctx context.Context, url string, data interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("got response code: %v", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if err := json.Unmarshal(body, data); err != nil {
		return fmt.Errorf("unmarshaling response body: %w", err)
	}
	return nil
}

// Weather represents weather conditions
// in a geographical region.
type Weather struct {
	Summary string
	Temp    float64
}

// String implements stringer interface.
func (w Weather) String() string {
	return fmt.Sprintf("%s %.1f°C", cases.Title(language.English).String(w.Summary), w.Temp)
}

// GetWeather returns current weather for given
// place and country using default client.
func GetWeather(location string) (Weather, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c, err := NewClient()
	if err != nil {
		return Weather{}, err
	}
	return c.GetWeather(ctx, location)
}

// RunCLI is a main function that runs the cli machinery.
func RunWeatherCLI() int {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s LOCATION\n\nExmple: %[1]s London,UK\n", os.Args[0])
		return 1
	}
	location := strings.Join(os.Args[1:], " ")
	w, err := GetWeather(location)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Fprintln(os.Stdout, w)
	return 0
}

func RunForecastCLI() int {
	return 0
}
