package cheek

import (
	"os"
	"strings"
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
	logger := NewLogger("debug", nil, b, os.Stdout)

	go func() {
		err := RunSchedule(logger, Config{DBPath: "tmpdb.sqlite3"}, "../testdata/jobs1.yaml")
		if err != nil {
			panic(err)
		}
	}()

	time.Sleep(6 * time.Second)
	spew.Dump(b.String())
	if err := proc.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}

	time.Sleep(6 * time.Second)
	assert.Contains(t, b.String(), "Job triggered")
	assert.Contains(t, b.String(), "Shutting down scheduler due to context cancellation")

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

func TestDisableConcurrentExecution(t *testing.T) {
	// Test that when disable_concurrent_execution is true, only one instance runs
	// and when false, multiple instances can run concurrently

	// Create a test schedule with jobs that take 3 seconds but run every second
	testScheduleYAML := `
jobs:
  concurrent_disabled:
    command: 
      - sh
      - -c
      - "echo 'start_non_concurrent'; sleep 10; echo 'end_non_concurrent'"
    cron: "* * * * * *"  # every second
    disable_concurrent_execution: true
  concurrent_enabled:
    command: 
      - sh
      - -c
      - "echo 'start_concurrent'; sleep 10; echo 'end_concurrent'"
    cron: "* * * * * *"  # every second
    disable_concurrent_execution: false
`

	// Write test schedule to temp file
	tmpFile, err := os.CreateTemp("", "test_schedule_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	if _, err := tmpFile.WriteString(testScheduleYAML); err != nil {
		t.Fatal(err)
	}
	_ = tmpFile.Close()

	// Create buffer to capture logs
	b := new(tsBuffer)
	logger := NewLogger("debug", nil, b, os.Stdout)

	// Load the schedule
	s, err := loadSchedule(logger, Config{}, tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Create a done channel to signal scheduler completion
	done := make(chan struct{})

	// Run the scheduler with completion signaling
	go func() {
		defer close(done)
		s.Run()
	}()

	time.Sleep(5 * time.Second)

	// Send interrupt signal to stop scheduler
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}
	if err := proc.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}

	// Wait for scheduler to fully shut down
	<-done

	// Collect job logs from individual jobs
	var allJobLogs strings.Builder

	// Access job runs that are stored in memory (without DB dependency)
	for _, job := range s.Jobs {

		for _, run := range job.Runs {
			spew.Dump(run)
			allJobLogs.WriteString(run.Log)
		}
	}

	completeLogOutput := allJobLogs.String()

	// Count occurrences of start messages
	nonConcurrentStarts := strings.Count(completeLogOutput, "start_non_concurrent")
	concurrentStarts := strings.Count(completeLogOutput, "start_concurrent")

	// With disable_concurrent_execution: true, we should see exactly 1 start
	// because subsequent executions are blocked while the first 3-second job is running
	assert.Equal(t, 1, nonConcurrentStarts, "Expected exactly 1 start for non-concurrent job")

	// With disable_concurrent_execution: false, we should see more than 1 start
	// because jobs can overlap (8 seconds runtime with 3-second jobs starting every second)
	assert.Greater(t, concurrentStarts, 1, "Expected more than 1 start for concurrent job")
}
