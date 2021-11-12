package meteo_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/meteo"
	"github.com/qba73/meteo/geonames"
	"github.com/qba73/meteo/geonames/testhelper"
)

func TestCreateNewMeteoClient(t *testing.T) {
	t.Parallel()

	resolver, err := geonames.NewClient("DummyUser")
	if err != nil {
		t.Fatal(err)
	}
	c, err := meteo.NewYrClient(resolver)
	if err != nil {
		t.Fatal(err)
	}
	_ = c
}

func TestCreateNewNorwayClientWithInvalidUserAgent(t *testing.T) {
	t.Parallel()

	resolver, err := geonames.NewClient("DummyUser")
	if err != nil {
		t.Fatal(err)
	}
	_, err = meteo.NewYrClient(resolver, meteo.WithUserAgent(""))
	if err == nil {
		t.Errorf("invalid user agent string should return error")
	}
}

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

	mux.HandleFunc("/wikipediaSearchJSON", func(rw http.ResponseWriter, r *http.Request) {
		testFile := "testdata/response-geoname-wikipedia-single.json"
		wantURI := "/wikipediaSearchJSON?maxRows=1&q=Castlebar&title=Castlebar&countryCode=IE&username=DummyUser"
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

	resolver, err := geonames.NewClient(
		"DummyUser",
		geonames.WithBaseURL(ts.URL),
	)
	if err != nil {
		t.Fatal(err)
	}

	client, err := meteo.NewYrClient(
		resolver,
		meteo.WithBaseURL(ts.URL),
	)
	if err != nil {
		t.Fatal(err)
	}

	got, err := client.GetForecast("Castlebar", "IE")
	if err != nil {
		t.Errorf("error getting forecast data, %v", err)
	}
	want := meteo.Weather{
		Summary: "rain",
		Temp:    13.7,
	}

	if !cmp.Equal(fmt.Sprint(want), got) {
		t.Error(cmp.Diff(fmt.Sprint(want), got))
	}
}
