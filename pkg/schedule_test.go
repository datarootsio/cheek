package cheek

import (
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestScheduleRun(t *testing.T) {
	// rough test
	// just tries to see if we can get to a job trigger
	// and to see that exit signals are received correctly
	viper.Set("port", "9999")
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	b := new(tsBuffer)
	logger := NewLogger(false, "debug", b)
	go func() {
		RunSchedule(logger, Config{}, "../testdata/jobs1.yaml")
	}()

	time.Sleep(3 * time.Second)
	if err := proc.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)
	spew.Dump(b.String())
	assert.Contains(t, b.String(), "Job triggered")
	// assert.Contains(t, b.String(), "interrupt signal received")

	// check that job gets triggered by other job
	assert.Contains(t, b.String(), "\"trigger\":\"job[foo]")
}
