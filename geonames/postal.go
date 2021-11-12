package geonames

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
)

const (
	postalBasePath = "postalCodeSearchJSON"
)

type PostalCodes struct {
	Codes []PostalCode `json:"postalCodes"`
}

type PostalCode struct {
	PlaceName   string  `json:"placeName"`
	AdminName1  string  `json:"adminName1"`
	AdminName2  string  `json:"adminName2"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	CountryCode string  `json:"countryCode"`
	PostalCode  string  `json:"postalCode"`
	AdminCode1  string  `json:"adminCode1"`
	AdminCode2  string  `json:"adminCode2"`
}

type PostalCodesService struct {
	cl *client
}

// Resolve knows how to retrieve lat and long coordinates
// for the given place name and country.
func (ps PostalCodesService) Get(placeName, countryCode string) (PostalCodes, error) {
	if placeName == "" || countryCode == "" {
		return PostalCodes{}, errors.New("nil place name or countryCode")
	}
	u, err := ps.makePostalURL(placeName, countryCode)
	if err != nil {
		return PostalCodes{}, err
	}

	req, err := prepareGetRequest(u)
	if err != nil {
		return PostalCodes{}, err
	}
	res, err := ps.cl.HTTPClient.Do(req)
	if err != nil {
		return PostalCodes{}, err
	}
	defer res.Body.Close()

	var pc PostalCodes
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PostalCodes{}, fmt.Errorf("reading response body %w", err)
	}

	if err := json.Unmarshal(data, &pc); err != nil {
		return PostalCodes{}, fmt.Errorf("unmarshalling data, %w", err)
	}
	return pc, nil
}

func (ps PostalCodesService) makePostalURL(placeName, countryCode string) (string, error) {
	prms := url.Values{
		"placename": {placeName},
		"country":   {countryCode},
		"username":  {ps.cl.UserName},
	}
	basePostal := fmt.Sprintf("%s/%s", ps.cl.BaseURL, postalBasePath)
	return makeURL(basePostal, prms)
}
