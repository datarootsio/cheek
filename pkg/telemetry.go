package cheek

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type ET struct{}

func (et ET) PhoneHome() ([]byte, error) {
	if !viper.IsSet("noTelemetry") {
		return []byte{}, errors.New("noTelemetry default not set")
	}
	if viper.GetBool("noTelemetry") {
		return []byte{}, nil
	}

	if !viper.IsSet("phoneHomeURL") {
		return []byte{}, errors.New("phoneHomeURL not set")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", viper.GetString("phoneHomeURL"), nil)
	if err != nil {
		return []byte{}, err
	}

	q := req.URL.Query()
	q.Add("version", Version)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
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
