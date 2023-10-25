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
	var err error
	var resp_body []byte

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		// mirror this
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(body))
		t.Log(string(body))
	}))

	defer testServer.Close()

	// test generic webhook
	jr := JobRun{
		Status:      0,
		Name:        "test",
		TriggeredBy: "cron",
		Log:         "this is a random log statement\nwith multiple lines\nand stuff",
	}

	resp_body, err = JobRunWebhookCall(&jr, testServer.URL, "generic")
	assert.NoError(t, err)

	jr2 := JobRun{}
	err = json.NewDecoder(bytes.NewBuffer(resp_body)).Decode(&jr2)
	assert.NoError(t, err)

	assert.Equal(t, jr, jr2)

	// test slack webhook
	resp_body, err = JobRunWebhookCall(&jr, testServer.URL, "slack")
	assert.NoError(t, err)
	assert.Contains(t, string(resp_body), "text\":\"test (exitcode 0)")

}
