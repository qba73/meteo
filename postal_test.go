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

func TestCreateNewGeoNamesClient(t *testing.T) {
	t.Parallel()
	var c *meteo.GeoNamesClient
	c, err := meteo.NewGeoNamesClient("user")
	if err != nil {
		t.Fatal(err)
	}
	_ = c
}

func TestCreateNewGeoNamesClientWithoutUserName(t *testing.T) {
	t.Parallel()
	_, err := meteo.NewGeoNamesClient("")
	if err == nil {
		t.Fatal("create client without user should return err")
	}
}

func TestCreateNewGeoNamesClientWithUser(t *testing.T) {
	t.Parallel()
	c, err := meteo.NewGeoNamesClient("User")
	if err != nil {
		t.Fatal(err)
	}
	want := "User"
	if want != c.UserName {
		t.Errorf("want %s, got %s", want, c.UserName)
	}
}

func TestGetCoordinatesSingleGeoNames(t *testing.T) {
	t.Parallel()

	testFile := "testdata/response-geoname-single.json"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wantReqURL := "/postalCodeSearchJSON?country=IE&placename=Castlebar&username=UserName"
		gotReqURL := r.RequestURI
		if wantReqURL != gotReqURL {
			t.Errorf("want %q for wikipedia URL, got %q", wantReqURL, gotReqURL)
		}

		f, err := os.Open(testFile)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		io.Copy(w, f)
	}))
	defer ts.Close()

	c, err := meteo.NewGeoNamesClient("UserName")
	if err != nil {
		t.Fatal(err)
	}
	c.BaseURL = ts.URL

	got, err := c.GetCoordinates("Castlebar", "IE")
	if err != nil {
		t.Fatalf("GetCoordinates(\"Castlebar\", \"IE\") got err %v", err)
	}
	want := meteo.Place{
		Lng:         -9.3,
		Lat:         53.85,
		PlaceName:   "Castlebar",
		CountryCode: "IE",
	}

	if !cmp.Equal(want, got) {
		t.Errorf("GetCoordinates('Castlebar', 'IE') \n%s", cmp.Diff(want, got))
	}
}
