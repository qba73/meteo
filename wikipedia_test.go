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

func TestNewWikipediaClient(t *testing.T) {
	t.Parallel()
	var err error
	var c *meteo.WikipediaClient
	c, err = meteo.NewWikipediaClient("UserName")
	if err != nil {
		t.Fatal(err)
	}
	_ = c
}

func TestNewWikipediaClientWithoutUserName(t *testing.T) {
	t.Parallel()
	_, err := meteo.NewWikipediaClient("")
	if err == nil {
		t.Fatal(err)
	}
}

func TestNewWikipediaClientWithUserName(t *testing.T) {
	t.Parallel()
	got, err := meteo.NewWikipediaClient("UserName")
	if err != nil {
		t.Fatal(err)
	}
	want := "UserName"
	if want != got.UserName {
		t.Errorf("want %s, got %s", want, got.UserName)
	}
}

func TestGetCoordinatesSingleGeoName(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/response-geoname-wikipedia.json")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		io.Copy(rw, f)
	}))
	defer ts.Close()

	c, err := meteo.NewWikipediaClient("UserName")
	if err != nil {
		t.Fatal(err)
	}
	c.BaseURL = ts.URL

	got, err := c.GetCoordinates("Castlebar", "IE")
	if err != nil {
		t.Fatalf("GetCoordinates('Castlebar', 'IE') got err %v", err)
	}
	want := meteo.Place{
		Lng:         -9.2988,
		Lat:         53.8608,
		PlaceName:   "Castlebar",
		CountryCode: "IE",
	}

	if !cmp.Equal(want, got) {
		t.Errorf("GetCoordinates('Castlebar', 'IE') \n%s", cmp.Diff(want, got))
	}
}
