package butt

import (
	"bytes"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestThatTuiStarts(t *testing.T) {
	// very rough test
	b := bytes.Buffer{}
	hp := "9999"
	go func() {
		RunSchedule("../testdata/jobs1.yaml", true, hp, true, "debug")
	}()
	time.Sleep(500 * time.Millisecond)
	log.Logger = zerolog.New(&b)

	go func() {
		TUI(hp)
	}()

	time.Sleep(1 * time.Second)
	assert.Contains(t, b.String(), "Better Unified")

}
