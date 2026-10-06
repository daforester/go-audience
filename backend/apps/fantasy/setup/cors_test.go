package setup

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestAllowedOriginsFromEnvironment(t *testing.T) {
	viper.Reset()
	viper.SetEnvPrefix("FANTASY")
	viper.AutomaticEnv()
	Defaults()

	t.Setenv("FANTASY_CORS_ALLOW_ORIGINS", "https://ffc.furrific.uk  https://admin.furrific.uk/")

	expected := []string{"https://ffc.furrific.uk", "https://admin.furrific.uk"}
	if got := AllowedOrigins(); !reflect.DeepEqual(got, expected) {
		t.Fatalf("origins from env: %v", got)
	}
}

func TestOriginMatchesHost(t *testing.T) {
	cases := []struct {
		origin string
		host   string
		match  bool
	}{
		{"https://ffc-api.furrific.uk", "ffc-api.furrific.uk", true},
		{"http://192.168.1.20:5173", "192.168.1.20:5173", true},
		{"https://FFC-API.furrific.uk", "ffc-api.furrific.uk", true},
		{"https://ffc-api.furrific.uk", "ffc-api.furrific.uk:8443", false},
		{"https://evil.example", "ffc-api.furrific.uk", false},
		{"null", "ffc-api.furrific.uk", false},
		{"file://ffc-api.furrific.uk", "ffc-api.furrific.uk", false},
		{"https://ffc-api.furrific.uk", "", false},
	}
	for _, c := range cases {
		if got := OriginMatchesHost(c.origin, c.host); got != c.match {
			t.Errorf("OriginMatchesHost(%q, %q) = %v, want %v", c.origin, c.host, got, c.match)
		}
	}
}
