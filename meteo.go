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

// NewClient knows how to construct a new meteo client.
func NewClient(apikey string) *Client {
	return &Client{
		BaseURL: baseURL,
		apiKey:  apikey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

/*
type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type successResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
*/
