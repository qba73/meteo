package meteo_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/qba73/meteo"
)

func TestWeatherStringFormat(t *testing.T) {
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
