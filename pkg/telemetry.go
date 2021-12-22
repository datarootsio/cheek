package cheek

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ET struct{}

func (et ET) PhoneHome(phoneHomeUrl string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", phoneHomeUrl, nil)
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

	return resp_body, nil
}
