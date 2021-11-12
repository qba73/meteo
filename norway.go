package meteo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/qba73/meteo/geonames"
)

const (
	libVersion = "0.0.1"
	source     = "https://github.com/qba73/meteo"
	userAgent  = "Meteo/" + libVersion + " " + source
	baseURL    = "https://api.met.no"
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

type NameResolver interface {
	Resolve(place, country string) (float64, float64, error)
}

type option func(*yrclient) error

// WithUserAgent is a func option to
// customise User-Agent header used internally
// by the yrclient.
func WithUserAgent(ua string) option {
	return func(nc *yrclient) error {
		if ua == "" {
			return errors.New("user agent not provided")
		}
		nc.userAgent = ua
		return nil
	}
}

// WithHTTPClient sets a new HTTP client for YR Client.
func WithHTTPClient(hc *http.Client) option {
	return func(c *yrclient) error {
		if hc == nil {
			return errors.New("nil http client provided")
		}
		c.httpClient = hc
		return nil
	}
}

// WithBaseURL sets a new URL for the Yr client.
// It errors if the provided base url is an empty string.
func WithBaseURL(u string) option {
	return func(c *yrclient) error {
		if u == "" {
			return errors.New("nil base URL")
		}
		c.baseURL = u
		return nil
	}
}

// yrclient represents a weather client
// for the Norwegian Meteorological Institute.
type yrclient struct {
	resolver   NameResolver
	userAgent  string
	baseURL    string
	httpClient *http.Client
}

// NewNorwayClient knows how to construct a new client.
func NewYrClient(resolver NameResolver, opts ...option) (*yrclient, error) {
	c := yrclient{
		userAgent: userAgent,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
		resolver: resolver,
	}
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}
	return &c, nil
}

// GetForecast takes place and country code and
// returns weather summary and air temperature.
func (c yrclient) GetForecast(place, country string) (string, error) {
	lat, lng, err := c.resolver.Resolve(place, country)
	if err != nil {
		return "", err
	}
	return c.getForecast(lat, lng)
}

func (c yrclient) getForecast(lat, lon float64) (string, error) {
	u, err := c.makeURL(lat, lon)
	if err != nil {
		return "", err
	}
	req, err := prepareRequest(u)
	if err != nil {
		return "", err
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body, %v", err)
	}

	var nf norwayForecast
	if err := json.Unmarshal(data, &nf); err != nil {
		return "", fmt.Errorf("unmarshalling data, %v", err)
	}

	if len(nf.Properties.Timeseries) < 1 {
		return "", fmt.Errorf("invalid response %+v", nf)
	}

	w := Weather{
		Summary: nf.Properties.Timeseries[0].Data.Next1Hours.Summary.SymbolCode,
		Temp:    nf.Properties.Timeseries[0].Data.Instant.Details.AirTemperature,
	}
	return fmt.Sprint(w), nil
}

func (c yrclient) makeURL(lat, lon float64) (string, error) {
	base, err := url.Parse(c.baseURL + "/weatherapi/locationforecast/2.0/compact")
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%.2f", lat))
	params.Add("lon", fmt.Sprintf("%.2f", lon))
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func prepareRequest(u string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", userAgent)
	return req, nil
}

// GetWeather returns current weather for given
// place and country using default client for the Norwegian
// meteorological Institute.
func GetWeather(place, country, username string) (string, error) {
	resolver, err := geonames.NewClient(
		username,
		geonames.WithUserAgent(userAgent),
	)
	if err != nil {
		return "", err
	}
	c, err := NewYrClient(resolver, WithUserAgent(userAgent))
	if err != nil {
		return "", err
	}
	return c.GetForecast(place, country)
}
