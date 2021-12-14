package cheek

import (
	"bytes"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestScheduleRun(t *testing.T) {
	// rough test
	// just tries to see if we can get to a job trigger
	b := bytes.Buffer{}

	go func() {
		RunSchedule("../testdata/jobs1.yaml", true, "9999", true, "debug")
	}()
	time.Sleep(100 * time.Millisecond)
	log.Logger = zerolog.New(&b)
	time.Sleep(2 * time.Second)
	assert.Contains(t, b.String(), "Job triggered")
}
