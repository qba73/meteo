package meteo_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/meteo"
)

func TestCreateNewMeteoClient(t *testing.T) {
	t.Parallel()
	var c *meteo.NorwayClient
	c, err := meteo.NewNorwayClient()
	if err != nil {
		t.Fatal(err)
	}
	_ = c
}

func TestCreateNewMeteoClientWithCustomUserAgent(t *testing.T) {
	t.Parallel()
	c, err := meteo.NewNorwayClient(
		meteo.WithUserAgent("CustomClient/1.0 https://customclient.com"),
	)
	if err != nil {
		t.Fatalf("creating client with custom agent, %s\n", err)
	}
	want := "CustomClient/1.0 https://customclient.com"
	got := c.UA
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestCreateNewNorwayClientWithInvalidUserAgent(t *testing.T) {
	t.Parallel()
	_, err := meteo.NewNorwayClient(meteo.WithUserAgent(""))
	if err == nil {
		t.Errorf("invalid user agent string should return error")
	}
}

func TestGetForecast(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/response-compact.json")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		io.Copy(rw, f)
	}))
	defer ts.Close()

	client, err := meteo.NewNorwayClient()
	if err != nil {
		t.Fatal(err)
	}
	client.BaseURL = ts.URL

	got, err := client.GetForecast(53.3, -6.2)
	if err != nil {
		t.Errorf("error getting forecast data, %v", err)
	}
	want := meteo.Weather{
		Summary: "rain",
		Temp:    13.7,
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
