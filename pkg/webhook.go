package cheek

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func JobRunWebhookCall(jr *JobRun, webhookURL string) ([]byte, error) {
	payload := bytes.Buffer{}
	if err := json.NewEncoder(&payload).Encode(jr); err != nil {
		return []byte{}, err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload.Bytes()))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	log.Debug().Str("job", jr.Name).Str("webhook_call", "response").Str("webhook_url", webhookURL).Msg(string(resp_body))

	return resp_body, nil
}
