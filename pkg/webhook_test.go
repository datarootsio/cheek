package cheek

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobRunWebhookCall(t *testing.T) {
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

	jr := JobRun{
		Status:      0,
		Name:        "test",
		TriggeredBy: "cron",
	}

	err, resp_body := JobRunWebhookCall(&jr, testServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	jr2 := JobRun{}
	json.NewDecoder(bytes.NewBuffer(resp_body)).Decode(&jr2)

	assert.Equal(t, jr, jr2)
}
