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

	"github.com/qba73/geonames"
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

// Location represents geo coordinates.
type Location struct {
	Lat  float64
	Long float64
}

func resolve(location string) (Location, error) {
	resolver, err := geonames.NewClient(os.Getenv("GEONAMES_USER"))
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
		UserAgent: fmt.Sprintf("Meteo/%s https://github.com/qba73/meteo", libVersion),
		BaseURL:   "https://api.met.no",
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
		Resolve: resolve,
	}
	return &c
}

// GetWeather returns current weather for given place.
//
// Place string should have format: "<place-name>,<country-code>",
// for example: "London,UK", "Dublin,IE", "Paris,FR", "Warsaw,PL".
func (c YRclient) GetWeather(place string) (Weather, error) {
	location, err := c.Resolve(place)
	if err != nil {
		return Weather{}, err
	}
	return c.getForecast(location)
}

// GetWeatherForCoordinates returns current weather for a place
// with given coordinates (lat, long)
func (c YRclient) GetWeatherForCoordinates(lat, long float64) (Weather, error) {
	l := Location{
		Lat:  lat,
		Long: long,
	}
	return c.getForecast(l)

}

func (c YRclient) getForecast(location Location) (Weather, error) {
	u, err := c.makeURL(location.Lat, location.Long)
	if err != nil {
		return Weather{}, err
	}

	//var nf norwayForecast
	var nf forecastResponseCompact
	if err := c.get(u, &nf); err != nil {
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

func (c YRclient) get(url string, data interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

// GetWeather returns current weather for given
// place and country using default client.
func GetWeather(location string) (Weather, error) {
	return NewYRClient().GetWeather(location)
}
