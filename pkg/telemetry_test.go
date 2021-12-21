package cheek

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestPhoneHome(t *testing.T) {
	b := new(tsBuffer)
	ConfigLogger(false, "debug", b)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		// mirror this
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(body))
	}))

	defer testServer.Close()

	viper.Set("phoneHomeURL", testServer.URL)
	viper.Set("noTelemetry", false)

	et := ET{}
	_, err := et.PhoneHome()
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, b.String(), "phoned home")

	// respect to not phone home
	viper.Set("noTelemetry", true)
	b.Reset()

	_, err = et.PhoneHome()
	if err != nil {
		t.Fatal(err)
	}

	assert.NotContains(t, b.String(), "phoned home")
}
