package geonames_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/qba73/meteo/geonames"
	"github.com/qba73/meteo/geonames/testhelper"
)

func TestNewWikipediaClient(t *testing.T) {
	t.Parallel()
	c, err := geonames.NewClient("DummyUser")
	if err != nil {
		t.Fatal(err)
	}
	_ = c
}

func TestNewWikipediaClientWithoutUserName(t *testing.T) {
	t.Parallel()
	_, err := geonames.NewClient("")
	if err == nil {
		t.Fatal(err)
	}
}

func TestGetSingleGeoName(t *testing.T) {
	t.Parallel()

	testFile := "testdata/response-geoname-wikipedia-single.json"
	wantReqURI := "/wikipediaSearchJSON?q=Castlebar&title=Castlebar&countryCode=IE&maxRows=1&username=DummyUser"

	ts := testhelper.NewTestServer(testFile, wantReqURI, t)
	defer ts.Close()

	client, err := geonames.NewClient(
		"DummyUser",
		geonames.WithBaseURL(ts.URL),
	)
	if err != nil {
		t.Fatal(err)
	}

	place, counry := "Castlebar", "IE"
	resultLimit := 1

	got, err := client.Wikipedia.Get(place, counry, resultLimit)
	if err != nil {
		t.Fatalf("Get(%q, %q, %q) got err %v", place, counry, resultLimit, err)
	}
	want := geonames.WikiResponse{
		Geonames: []geonames.Geoname{
			{
				Summary:      "Castlebar is the county town of County Mayo, Ireland. It is in the middle of the county and is its largest town by population. A campus of Galway-Mayo Institute of Technology and the Country Life section of the National Museum of Ireland are two important local amenities (...)",
				Elevation:    41,
				GeoNameID:    2965654,
				Lat:          53.8608,
				Lng:          -9.2988,
				CountryCode:  "IE",
				Rank:         100,
				Lang:         "en",
				Title:        "Castlebar",
				WikipediaURL: "en.wikipedia.org/wiki/Castlebar",
			},
		},
	}

	if !cmp.Equal(want, got, cmpopts.IgnoreFields(geonames.Geoname{}, "Summary")) {
		t.Errorf(
			"Get(%q, %q, %q) got \n%s\n", place, counry, resultLimit,
			cmp.Diff(want, got),
		)
	}
}
