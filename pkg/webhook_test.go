package cheek

import (
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
		defer func() { _ = r.Body.Close() }()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		// mirror this
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, string(body))
		t.Log(string(body))
	}))

	defer testServer.Close()
	statusCode := 0
	// test generic webhook
	jr := JobRun{
		Status:      &statusCode,
		Name:        "test",
		TriggeredBy: "cron",
		Log:         "this is a random log statement\nwith multiple lines\nand stuff",
	}

	var wh webhook
	wh = NewDefaultWebhook(testServer.URL)
	resp_body, err = wh.Call(&jr)
	assert.NoError(t, err)
	assert.Contains(t, string(resp_body), `{"status":0,"log":"this is a random log statement\nwith multiple lines\nand stuff","name":"test","triggered_at":"0001-01-01T00:00:00Z","triggered_by":"cron"}`)

	// test slack webhook
	wh = NewSlackWebhook(testServer.URL)
	resp_body, err = wh.Call(&jr)
	assert.NoError(t, err)
	assert.Contains(t, string(resp_body), `{"text":"test (exitcode 0):\nthis is a random log statement\nwith multiple lines\nand stuff"}`)

	// test discord webhook
	wh = NewDiscordWebhook(testServer.URL)
	resp_body, err = wh.Call(&jr)
	assert.NoError(t, err)
	assert.Contains(t, string(resp_body), `{"content":"test (exitcode 0):\nthis is a random log statement\nwith multiple lines\nand stuff"}`)
}
