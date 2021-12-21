package cheek

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadLogs(t *testing.T) {
	j := &JobSpec{
		Cron:    "* * * * *",
		Name:    "test",
		Command: []string{"echo", "bar"},
	}

	j.execCommand("test", true)

	// log loading goes on job name basis
	// let's recreate
	j = &JobSpec{
		Name: "test",
	}

	j.loadRuns()
	assert.Greater(t, len(j.runs), 0)
}

func TestJobRun(t *testing.T) {
	j := &JobSpec{
		Cron:    "* * * * *",
		Name:    "test",
		Command: []string{"echo", "bar"},
	}

	jr := j.execCommand("test", true)
	assert.Equal(t, jr.Status, 0)
}

func TestJobRunNoCommand(t *testing.T) {
	j := &JobSpec{
		Cron: "* * * * *",
		Name: "test",
	}

	jr := j.execCommand("test", true)
	assert.NotEqual(t, jr.Status, 0)
}

func TestJobNonZero(t *testing.T) {
	j := &JobSpec{
		Cron: "* * * * *",
		Name: "test",
		Command: []string{
			"la", "--moo",
		},
	}

	jr := j.execCommand("test", true)
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
	s.Jobs["Bertha"].OnSuccess.TriggerJobs = []string{"IDontExist"}

	assert.Error(t, s.Validate())
}
