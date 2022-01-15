package cheek

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhoneHome(t *testing.T) {
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

	et := ET{}
	_, err := et.PhoneHome(testServer.URL)
	if err != nil {
		// no error means success
		t.Fatal(err)
	}

	_, err = et.PhoneHome("http://non-existtant-host")
	assert.Error(t, err)

	_, err = et.PhoneHome("httx://non-existtant-protocol")
	assert.Error(t, err)
}
