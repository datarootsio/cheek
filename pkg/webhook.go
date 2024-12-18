package cheek

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type webhook interface {
	Call(jr *JobRun) ([]byte, error)
	URL() string
	Name() string
}

// Discord Webhook

type discordWebhook struct {
	endpoint string
}

func NewDiscordWebhook(endpoint string) discordWebhook {
	return discordWebhook{endpoint}
}

func (dw discordWebhook) Call(jr *JobRun) ([]byte, error) {
	type discordPayload struct {
		Content string `json:"content"`
	}
	payload := bytes.Buffer{}
	msg := fmt.Sprintf("%s (exitcode %v):\n%s", jr.Name, *jr.Status, jr.Log)
	d := discordPayload{
		Content: msg[:min(len(msg), 2000)], // discord accepts a max. of 2000 chars
	}
	if err := json.NewEncoder(&payload).Encode(d); err != nil {
		return nil, err
	}

	resp, err := http.Post(dw.endpoint, "application/json", bytes.NewBuffer(payload.Bytes()))
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

func (dw discordWebhook) URL() string {
	return dw.endpoint
}

func (dw discordWebhook) Name() string {
	return "discord"
}

// Slack Webhook

type slackWebhook struct {
	endpoint string
}

func NewSlackWebhook(endpoint string) slackWebhook {
	return slackWebhook{endpoint}
}

func (dw slackWebhook) Call(jr *JobRun) ([]byte, error) {
	type slackPayload struct {
		Text string `json:"text"`
	}
	payload := bytes.Buffer{}
	msg := fmt.Sprintf("%s (exitcode %v):\n%s", jr.Name, *jr.Status, jr.Log)
	d := slackPayload{
		Text: msg[:min(len(msg), 40000)], // slack accepts a max. of 40000 chars
	}
	if err := json.NewEncoder(&payload).Encode(d); err != nil {
		return nil, err
	}

	resp, err := http.Post(dw.endpoint, "application/json", bytes.NewBuffer(payload.Bytes()))
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

func (dw slackWebhook) URL() string {
	return dw.endpoint
}

func (dw slackWebhook) Name() string {
	return "slack"
}

// Default Webhook

type defaultWebhook struct {
	endpoint string
}

func NewDefaultWebhook(endpoint string) defaultWebhook {
	return defaultWebhook{endpoint}
}

func (dw defaultWebhook) Call(jr *JobRun) ([]byte, error) {
	payload := bytes.Buffer{}
	if err := json.NewEncoder(&payload).Encode(jr); err != nil {
		return []byte{}, err
	}

	resp, err := http.Post(dw.endpoint, "application/json", bytes.NewBuffer(payload.Bytes()))
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

func (dw defaultWebhook) URL() string {
	return dw.endpoint
}

func (dw defaultWebhook) Name() string {
	return "generic"
}
