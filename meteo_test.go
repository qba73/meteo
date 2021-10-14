package meteo_test

import (
	"testing"

	"github.com/qba73/meteo"
)

func TestNewClient(t *testing.T) {
	meteo.NewClient("APIKEY")
}
