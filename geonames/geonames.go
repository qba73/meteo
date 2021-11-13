package geonames

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	libraryVersion = "0.1"
	userAgent      = "geonames/" + libraryVersion
	baseURL        = "http://api.geonames.org"
	mediaType      = "application/json"
)

// Resolver interface is used by an Meteo Client
// to obtain geo coordinates for given place located in
// a country identified by country id.
type Resolver interface {
	// Resolve takes place and country code
	// and returns lat lng coordinates.
	Resolve(place string) (string, error)
}

type Option func(*client) error

// WithUserAgent knows how to add custom user agent to web requests.
func WithUserAgent(ua string) Option {
	return func(c *client) error {
		if ua == "" {
			return errors.New("nil user agent")
		}
		c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
		return nil
	}
}

// WithBaseURL knows how to set base URL for the API client.
func WithBaseURL(bu string) Option {
	return func(c *client) error {
		c.BaseURL = bu
		return nil
	}
}

// WithHeaders knows how to set custom HTTP headers for each API request.
func WithHeaders(headers map[string]string) Option {
	return func(c *client) error {
		for k, v := range headers {
			c.headers[k] = v
		}
		return nil
	}
}

// WithHTTPClient knows how to set a custom HTTP client used
// by the GeoNames client.
func WithHTTPClient(cl *http.Client) Option {
	return func(c *client) error {
		c.HTTPClient = cl
		return nil
	}
}

// Client is a client used for communicating
// with GeoNames web service.
type client struct {
	// UserName is a user name chosen when registered for GeoNames.org
	UserName string

	UserAgent  string
	BaseURL    string
	HTTPClient *http.Client

	// Optional HTTP headers to set for each API request.
	headers map[string]string

	Wikipedia   WikipediaService
	PostalCodes PostalCodesService
}

// NewClient knows how to create a client for GeoNames Web service.
// It returns an error when user name is not provided. Note that
// providing user name that is not registered on GeoNames.org
// will result with HTTP errors 403.
func NewClient(userName string, opts ...Option) (*client, error) {
	if userName == "" {
		return nil, errors.New("missing user name")
	}
	c := client{
		UserAgent: userAgent,
		UserName:  userName,
		BaseURL:   baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
	for _, opt := range opts {
		if err := opt(&c); err != nil {
			return nil, err
		}
	}
	c.Wikipedia = WikipediaService{cl: &c}
	c.PostalCodes = PostalCodesService{cl: &c}
	return &c, nil
}

func (c client) Resolve(place, country string) (float64, float64, error) {
	res, err := c.Wikipedia.Get(place, country, 1)
	if err != nil {
		return 0, 0, err
	}

	if len(res.Geonames) < 1 {
		return 0, 0, fmt.Errorf("place %s in country %s not found", place, country)
	}
	return res.Geonames[0].Lat, res.Geonames[0].Lng, nil
}

// makeURL knows how to create encoded url with provided query parameters.
func makeURL(base string, params url.Values) (string, error) {
	b, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("parsing base url, %v", err)
	}
	b.RawQuery = params.Encode()
	return b.String(), nil
}

func prepareGetRequest(u string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("User-Agent", userAgent)
	return req, nil
}
