package meteo_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/meteo"
)

func newTestServerWithPathValidator(testFile string, wantURI string, t *testing.T) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		gotReqURI := r.RequestURI
		verifyURIs(wantURI, gotReqURI, t)

		f, err := os.Open(testFile)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		_, err = io.Copy(rw, f)
		if err != nil {
			t.Fatalf("copying data from file %s to test HTTP server: %v", testFile, err)
		}
	}))
	return ts
}

// verifyURIs is a test helper function that verifies if provided URIs are equal.
func verifyURIs(wanturi, goturi string, t *testing.T) {
	wantU, err := url.Parse(wanturi)
	if err != nil {
		t.Fatalf("error parsing URL %q, %v", wanturi, err)
	}
	gotU, err := url.Parse(goturi)
	if err != nil {
		t.Fatalf("error parsing URL %q, %v", wanturi, err)
	}

	if !cmp.Equal(wantU.Path, gotU.Path) {
		t.Fatalf(cmp.Diff(wantU.Path, gotU.Path))
	}

	wantQuery, err := url.ParseQuery(wantU.RawQuery)
	if err != nil {
		t.Fatal(err)
	}
	gotQuery, err := url.ParseQuery(gotU.RawQuery)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(wantQuery, gotQuery) {
		t.Fatalf("URIs are not equal, \n%s", cmp.Diff(wantQuery, gotQuery))
	}
}

func TestClientRequestsWeatherWithValidPathAndParams(t *testing.T) {
	t.Parallel()

	var called bool
	wantURI := "/weatherapi/locationforecast/2.0/compact?lat=53.86&lon=-9.30"
	testFile := "testdata/response-compact.json"

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		gotReqURI := r.RequestURI
		verifyURIs(wantURI, gotReqURI, t)

		f, err := os.Open(testFile)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		_, err = io.Copy(rw, f)
		if err != nil {
			t.Fatalf("copying data from file %s to test HTTP server: %v", testFile, err)
		}
		called = true
	}))
	defer ts.Close()

	client := meteo.NewClient()
	client.BaseURL = ts.URL
	client.Resolve = func(ctx context.Context, location string) (meteo.Location, error) {
		return meteo.Location{
			Lat:  53.86,
			Long: -9.30,
		}, nil
	}
	_, err := client.GetWeather(context.Background(), "Dublin,IE")
	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error("handler not called")
	}
}

func TestClientReadsCurrentWeatherOnValidInput(t *testing.T) {
	t.Parallel()

	testFile := "testdata/response-compact.json"
	wantURI := "/weatherapi/locationforecast/2.0/compact?lat=53.86&lon=-9.30"
	ts := newTestServerWithPathValidator(testFile, wantURI, t)
	defer ts.Close()

	client := meteo.NewClient()
	client.BaseURL = ts.URL
	client.Resolve = func(ctx context.Context, location string) (meteo.Location, error) {
		return meteo.Location{
			Lat:  53.86,
			Long: -9.30,
		}, nil
	}
	got, err := client.GetWeather(context.Background(), "Castlebar,IE")
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

func TestClientFormatsWeatherInfoOnValidInput(t *testing.T) {
	t.Parallel()
	w := meteo.Weather{
		Summary: "sunny",
		Temp:    -3.12,
	}
	out := bytes.Buffer{}
	fmt.Fprint(&out, w)
	got := out.String()
	want := "Sunny -3.1°C"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}
