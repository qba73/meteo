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
	baseURLGeoNames = "http://api.geonames.org"
)

// Place represents response from GeoNames web service.
type Place struct {
	Lat         float64
	Lng         float64
	CountryCode string
	PlaceName   string
}

type coordinates struct {
	PostalCodes []struct {
		Lat         float64 `json:"lat"`
		Lng         float64 `json:"lng"`
		CountryCode string  `json:"countryCode"`
		PlaceName   string  `json:"placeName"`
	} `json:"postalCodes"`
}

type geoNameOption func(*GeoNamesClient) error

// WithGeoNamesUserAgent knows how to add custom user agent
// to web requests.
func WithGeoNamesUserAgent(ua string) geoNameOption {
	return func(gnc *GeoNamesClient) error {
		if ua == "" {
			return errors.New("nil user agent")
		}
		gnc.UA = ua
		return nil
	}
}

// GeoNamesClient is a client used for communicating
// with geo name web services.
type GeoNamesClient struct {
	UA         string
	UserName   string
	BaseURL    string
	HTTPClient *http.Client
}

// NewGeoNamesClient knows how to create a client for GeoNames Web services.
func NewGeoNamesClient(username string, opts ...geoNameOption) (*GeoNamesClient, error) {
	if username == "" {
		return nil, errors.New("missing user name")
	}

	c := GeoNamesClient{
		BaseURL:  baseURLGeoNames,
		UA:       userAgent,
		UserName: username,
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return &GeoNamesClient{}, err
		}
	}
	return &c, nil
}

// GetCoordinates knows how to retrieve lat and long coordinates
// for the given place name and country.
func (g GeoNamesClient) GetCoordinates(placeName, countryCode string) (Place, error) {
	u, err := g.makeURL(placeName, countryCode)
	if err != nil {
		return Place{}, err
	}
	req, err := prepareRequest(u)
	if err != nil {
		return Place{}, err
	}
	res, err := g.HTTPClient.Do(req)
	if err != nil {
		return Place{}, err
	}
	defer res.Body.Close()

	var co coordinates
	if err := json.NewDecoder(res.Body).Decode(&co); err != nil {
		return Place{}, fmt.Errorf("decoding response %w", err)
	}
	if len(co.PostalCodes) < 1 {
		return Place{}, fmt.Errorf("place %s in country %s not found", placeName, countryCode)
	}

	pc := Place{
		Lat:         co.PostalCodes[0].Lat,
		Lng:         co.PostalCodes[0].Lng,
		PlaceName:   co.PostalCodes[0].PlaceName,
		CountryCode: co.PostalCodes[0].CountryCode,
	}

	return pc, nil
}

func (g GeoNamesClient) makeURL(placeName, countryCode string) (string, error) {
	base, err := url.Parse(g.BaseURL + "/postalCodeSearchJSON")
	if err != nil {
		return "", fmt.Errorf("making url %w", err)
	}
	params := url.Values{}
	params.Add("placename", placeName)
	params.Add("country", countryCode)
	params.Add("username", g.UserName)
	base.RawQuery = params.Encode()
	return base.String(), nil
}

// GetCoordinates knows how to get Lat and Long coordinates for
// the given place and country using default geo client.
func GetCoordinates(placename, countryCode, username string) (Place, error) {
	c, err := NewGeoNamesClient(username)
	if err != nil {
		return Place{}, err
	}
	return c.GetCoordinates(placename, countryCode)
}
