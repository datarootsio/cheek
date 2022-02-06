package cheek

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type ET struct{}

func (et ET) PhoneHome(phoneHomeUrl string) ([]byte, error) {
	values := map[string]string{"version": Version}
	jsonBody, err := json.Marshal(values)
	if err != nil {
		return []byte{}, err
	}

	resp, err := http.Post(phoneHomeUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("cannot phone home, status %v", resp.StatusCode)
		return []byte{}, errors.New(msg)
	}

	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	log.Debug().Str("telemetry", "ET").Msg("ET phoned home")

	return resp_body, nil
}
