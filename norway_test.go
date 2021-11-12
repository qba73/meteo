package meteo_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/meteo"
	"github.com/qba73/meteo/geonames/testhelper"
)

func TestGetForecast(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()

	mux.HandleFunc("/weatherapi/locationforecast/2.0/compact", func(rw http.ResponseWriter, r *http.Request) {
		testFile := "testdata/response-compact.json"
		wantURI := "/weatherapi/locationforecast/2.0/compact?lat=53.86&lon=-9.30"
		gotURI := r.RequestURI

		testhelper.VerifyURIs(wantURI, gotURI, t)

		f, err := os.Open(testFile)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		_, err = io.Copy(rw, f)
		if err != nil {
			t.Fatalf("copying data from file %s to test HTTP server: %v", testFile, err)
		}
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	client := meteo.NewYRClient()
	client.BaseURL = ts.URL
	client.Resolve = func(location string) (meteo.Location, error) {
		return meteo.Location{
			Lat:  53.86,
			Long: -9.30,
		}, nil
	}
	got, err := client.GetForecast("Castlebar,IE")
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

func TestStringifyWeather(t *testing.T) {
	t.Parallel()
	input := meteo.Weather{
		Summary: "rain",
		Temp:    13.7,
	}
	want := "Rain 13.7°C"
	got := input.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
