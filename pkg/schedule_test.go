package cheek

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduleRun(t *testing.T) {
	// rough test
	// just tries to see if we can get to a job trigger
	// and to see that exit signals are received correctly

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	b := new(tsBuffer)
	lc := LogConfig{}
	lc.Init(false, "debug", b)
	defer lc.Close()
	go func() {
		RunSchedule("../testdata/jobs1.yaml", "9999", true)
	}()

	time.Sleep(2 * time.Second)
	proc.Signal(os.Interrupt)
	time.Sleep(1 * time.Second)
	assert.Contains(t, b.String(), "Job triggered")
	assert.Contains(t, b.String(), "interrupt signal received")
}
