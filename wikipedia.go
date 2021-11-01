package meteo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	maxRawResults = 10
)

type wikipediaPlaces struct {
	Geonames []struct {
		Lat         float64 `json:"lat"`
		Lng         float64 `json:"lng"`
		CountryCode string  `json:"countryCode"`
		Title       string  `json:"title"`
	} `json:"geonames"`
}

type wikipediaOption func(wk *WikipediaClient) error

func WikiWithUserAgent(ua string) wikipediaOption {
	return func(wk *WikipediaClient) error {
		wk.UA = ua
		return nil
	}
}

func WikiWithHTTPClient(hc *http.Client) wikipediaOption {
	return func(wk *WikipediaClient) error {
		wk.HTTPClient = hc
		return nil
	}
}

// WikipediaClient implements NameResolver interface.
type WikipediaClient struct {
	UA         string
	UserName   string
	BaseURL    string
	HTTPClient *http.Client
}

// NewWikipediaClient knows how to create a new client.
// The client impements NameResolver interface.
func NewWikipediaClient(username string) (*WikipediaClient, error) {
	if username == "" {
		return nil, errors.New("nil username")
	}
	c := WikipediaClient{
		UA:       userAgent,
		UserName: username,
		BaseURL:  baseURLGeoNames,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	return &c, nil
}

// GetCoordinates knows how to retrive geo coordinates for
// the given place name and country code.
func (w WikipediaClient) GetCoordinates(placeName, country string) (Place, error) {
	u, err := w.makeURL(placeName)
	if err != nil {
		return Place{}, err
	}
	req, err := prepareRequest(u)
	if err != nil {
		return Place{}, err
	}
	res, err := w.HTTPClient.Do(req)
	if err != nil {
		return Place{}, err
	}
	defer res.Body.Close()

	var wp wikipediaPlaces
	if err := json.NewDecoder(res.Body).Decode(&wp); err != nil {
		return Place{}, fmt.Errorf("decoding response %w", err)

	}
	if len(wp.Geonames) < 1 {
		return Place{}, fmt.Errorf("place %s in country %s not found", placeName, country)
	}
	return lookupPlace(wp, placeName, country), nil
}

func (w WikipediaClient) makeURL(placeName string) (string, error) {
	base, err := url.Parse(w.BaseURL + "/wikipediaSearchJSON")
	if err != nil {
		return "", err
	}
	params := url.Values{}
	params.Add("q", placeName)
	params.Add("maxRows", strconv.Itoa(maxRawResults))
	params.Add("username", w.UserName)
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func lookupPlace(places wikipediaPlaces, placeName, country string) Place {
	for _, pl := range places.Geonames {
		if pl.Title == placeName && pl.CountryCode == country {
			return Place{
				Lat:         pl.Lat,
				Lng:         pl.Lng,
				CountryCode: pl.CountryCode,
				PlaceName:   pl.Title,
			}
		}
	}
	return Place{}
}
