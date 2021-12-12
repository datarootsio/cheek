package butt

import (
	"bytes"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestLoadLogs(t *testing.T) {
	j := &JobSpec{
		Cron:    "* * * * *",
		Name:    "test",
		Command: []string{"echo", "bar"},
	}

	j.ExecCommand("test", true)

	// log loading goes on job name basis
	// let's recreate
	j = &JobSpec{
		Name: "test",
	}

	j.LoadRuns()
	assert.Greater(t, len(j.runs), 0)
}
func TestJobRun(t *testing.T) {
	j := &JobSpec{
		Cron:    "* * * * *",
		Name:    "test",
		Command: []string{"echo", "bar"},
	}

	jr := j.ExecCommand("test", true)
	assert.Equal(t, jr.Status, 0)
}

func TestJobRunNoCommand(t *testing.T) {
	j := &JobSpec{
		Cron: "* * * * *",
		Name: "test",
	}

	jr := j.ExecCommand("test", true)
	assert.NotEqual(t, jr.Status, 0)

}

func TestJobNonZero(t *testing.T) {
	j := &JobSpec{
		Cron: "* * * * *",
		Name: "test",
		Command: []string{
			"la", "--moo"},
	}

	jr := j.ExecCommand("test", true)
	assert.NotEqual(t, jr.Status, 0)

}

func TestJobRunInvalidSchedule(t *testing.T) {
	s := Schedule{}
	j := &JobSpec{
		Cron:    "MooIAmACow",
		Name:    "Bertha",
		Command: []string{"ls"},
	}
	s.Jobs = map[string]*JobSpec{}
	s.Jobs["Bertha"] = j

	assert.Error(t, s.Validate())
	// fix cron but add invalid ref
	s.Jobs["Bertha"].Cron = "* * * * *"
	s.Jobs["Bertha"].Triggers = []string{"IDontExist"}

	assert.Error(t, s.Validate())

}

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
