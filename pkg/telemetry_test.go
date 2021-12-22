package cheek

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
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
}
