package geonames

import (
	"net/url"
	"testing"

	"github.com/qba73/meteo/geonames/testhelper"
)

func TestCreateNewClient(t *testing.T) {
	t.Parallel()
	var c *client
	var err error
	c, err = NewClient("user")
	if err != nil {
		t.Fatal(err)
	}
	_ = c
}

func TestMakeURL(t *testing.T) {
	t.Parallel()

	baseURL := "http://geo.ie"
	params := url.Values{"title": []string{"Castlebar"}, "q": []string{"Castlebar"}}
	got, err := makeURL(baseURL, params)
	if err != nil {
		t.Fatal(err)
	}
	want := "http://geo.ie?q=Castlebar&title=Castlebar"
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestResolveGeoName(t *testing.T) {
	t.Parallel()

	testFile := "testdata/response-geoname-wikipedia-single.json"
	wantReqURI := "/wikipediaSearchJSON?q=Castlebar&title=Castlebar&countryCode=IE&maxRows=1&username=DummyUser"
	ts := testhelper.NewTestServer(testFile, wantReqURI, t)
	defer ts.Close()

	place, country := "Castlebar", "IE"
	client, err := NewClient(
		"DummyUser",
		WithBaseURL(ts.URL),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = client.Resolve(place, country)
	if err != nil {
		t.Error(err)
	}

}
