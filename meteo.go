package meteo

import (
	"net/http"
	"time"
)

const (
	baseURL = "https://"
)

// Client represents YR.no weather client.
type Client struct {
	apiKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient knows how to create Meteo service client.
func NewClient(apikey string) *Client {
	return &Client{
		BaseURL: baseURL,
		apiKey:  apikey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}
