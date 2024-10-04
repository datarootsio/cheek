package cheek

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type slackPayload struct {
	Text string `json:"text"`
}

func JobRunWebhookCall(jr *JobRun, webhookURL string, webhookType string) ([]byte, error) {
	payload := bytes.Buffer{}

	if webhookType == "slack" {
		d := slackPayload{
			Text: fmt.Sprintf("%s (exitcode %v):\n%s", jr.Name, *jr.Status, jr.Log),
		}

		if err := json.NewEncoder(&payload).Encode(d); err != nil {
			return []byte{}, err
		}

	} else {
		if err := json.NewEncoder(&payload).Encode(jr); err != nil {
			return []byte{}, err
		}
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
