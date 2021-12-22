package cheek

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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

	return resp_body, nil
}
