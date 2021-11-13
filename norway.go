package meteo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/qba73/meteo/geonames"
)

const (
	libVersion = "0.0.1"
	source     = "https://github.com/qba73/meteo"
)

type norwayForecast struct {
	Type       string     `json:"type"`
	Geometry   geometry   `json:"geometry"`
	Properties properties `json:"properties"`
}

type geometry struct {
	Type        string     `json:"type"`
	Coordinates []float64  `json:"coordinates"`
	Properties  properties `json:"properties"`
}

type properties struct {
	Meta       meta       `json:"meta"`
	Timeseries timeseries `json:"timeseries"`
}

type meta struct {
	UpdatedAt string            `json:"updated_at"`
	Units     map[string]string `json:"units"`
}

type timeseries []forecastEntry

type forecastEntry struct {
	Time string `json:"time"`
	Data data   `json:"data"`
}

type data struct {
	Instant struct {
		Details struct {
			AirTemperature float64 `json:"air_temperature"`
		}
	}
	Next1Hours struct {
		Summary struct {
			SymbolCode string `json:"symbol_code"`
		}
	} `json:"next_1_hours"`
}

type Location struct {
	Lat, Long float64
}

func resolve(location string) (Location, error) {
	uname := os.Getenv("GEO_USERNAME")
	resolver, err := geonames.NewClient(uname)
	if err != nil {
		return Location{}, err
	}
	place, country, err := toPlaceAndCountry(location)
	if err != nil {
		return Location{}, err
	}
	lat, long, err := resolver.Resolve(place, country)
	if err != nil {
		return Location{}, err
	}
	return Location{
		Lat:  lat,
		Long: long,
	}, nil
}

func toPlaceAndCountry(location string) (string, string, error) {
	bits := strings.Split(location, ",")
	if len(bits) < 2 {
		return "", "", fmt.Errorf("parsing place and country from location %s", location)
	}
	return bits[0], bits[1], nil
}

// YRclient represents a weather client
// for the Norwegian Meteorological Institute.
type YRclient struct {
	UserAgent  string
	BaseURL    string
	HTTPClient *http.Client
	Resolve    func(string) (Location, error)
}

// NewYRClient knows how to construct a new default client.
func NewYRClient() *YRclient {
	c := YRclient{
		UserAgent: "Meteo/" + libVersion + " " + source,
		BaseURL:   "https://api.met.no",
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
		Resolve: resolve,
	}
	return &c
}

// GetForecast takes place and returns weather
// summary and air temperature.
func (c YRclient) GetForecast(place string) (Weather, error) {
	location, err := c.Resolve(place)
	if err != nil {
		return Weather{}, err
	}
	return c.getForecast(location)
}

func (c YRclient) getForecast(location Location) (Weather, error) {
	u, err := c.makeURL(location.Lat, location.Long)
	if err != nil {
		return Weather{}, err
	}
	req, err := c.prepareRequest(u)
	if err != nil {
		return Weather{}, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return Weather{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Weather{}, fmt.Errorf("reading response body, %v", err)
	}

	var nf norwayForecast
	if err := json.Unmarshal(data, &nf); err != nil {
		return Weather{}, fmt.Errorf("unmarshalling data, %v", err)
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

func (c YRclient) makeURL(lat, lon float64) (string, error) {
	base, err := url.Parse(c.BaseURL + "/weatherapi/locationforecast/2.0/compact")
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%.2f", lat))
	params.Add("lon", fmt.Sprintf("%.2f", lon))
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (c YRclient) prepareRequest(u string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// GetWeather returns current weather for given
// place and country using default client.
func GetWeather(location string) (Weather, error) {
	return NewYRClient().GetForecast(location)
}
