package cheek

import (
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
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
	logger := NewLogger("debug", b, os.Stdout)

	go func() {
		err := RunSchedule(logger, Config{DBPath: "tmpdb.sqlite3"}, "../testdata/jobs1.yaml")
		if err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(61 * time.Second)
	spew.Dump(b.String())
	if err := proc.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	assert.Contains(t, b.String(), "Job triggered")
	assert.Contains(t, b.String(), "interrupt signal received")

	// check that job gets triggered by other job
	assert.Contains(t, b.String(), "\"trigger\":\"job[foo]")
}

func TestTZInfo(t *testing.T) {
	s := Schedule{
		Jobs:       map[string]*JobSpec{},
		TZLocation: "Africa/Bangui",
		log:        zerolog.Logger{},
		cfg:        NewConfig(),
	}
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}
	time1 := s.now()

	s = Schedule{
		Jobs:       map[string]*JobSpec{},
		TZLocation: "Europe/Amsterdam",
		log:        zerolog.Logger{},
		cfg:        NewConfig(),
	}
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	time2 := s.now()
	assert.NotEqual(t, time1.Sub(time2).Hours(), 0.0)
}
