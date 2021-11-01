package meteo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://api.met.no"
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

type option func(*NorwayClient) error

func WithUserAgent(ua string) option {
	return func(nc *NorwayClient) error {
		if ua == "" {
			return errors.New("user agent not provided")
		}
		nc.UA = ua
		return nil
	}
}

// NorwayClient represents a weather client
// for the Norwegian Meteorological Institute.
type NorwayClient struct {
	UA         string
	BaseURL    string
	HTTPClient *http.Client
}

func NewNorwayClient(opts ...option) (*NorwayClient, error) {
	c := NorwayClient{
		BaseURL: baseURL,
		UA:      userAgent,
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}

	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return &NorwayClient{}, err
		}
	}
	return &c, nil

}

func (c NorwayClient) GetForecast(lat, lon float64) (Weather, error) {
	u, err := c.makeURL(lat, lon)
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

	var nf norwayForecast
	if err := json.NewDecoder(res.Body).Decode(&nf); err != nil {
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

func (c NorwayClient) makeURL(lat, lon float64) (string, error) {
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

func (c NorwayClient) prepareRequest(u string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", userAgent)
	return req, nil
}

// GetWeather returns current weather for given
// Lat and Long using default client for the Norwegian
// meteorological Institute.
func GetWeather(lat, lon float64) (Weather, error) {
	c, err := NewNorwayClient()
	if err != nil {
		return Weather{}, err
	}
	return c.GetForecast(lat, lon)
}
