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
	c = meteo.NewNorwayClient()
	_ = c
}

func TestGetForecastCompact(t *testing.T) {
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

	client := meteo.NewNorwayClient()
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
