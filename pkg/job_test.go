package cheek

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLoadLogs(t *testing.T) {
	j := &JobSpec{
		Cron:    "* * * * *",
		Name:    "test",
		Command: []string{"echo", "bar"},
		cfg:     NewConfig(),
	}

	j.execCommand("test")

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
		cfg:     NewConfig(),
	}

	jr := j.execCommand("test")
	assert.Equal(t, jr.Status, 0)
}

func TestSpecialCron(t *testing.T) {
	j := &JobSpec{
		Cron:    "@10minutes",
		Name:    "test",
		Command: []string{"echo", "bar"},
		cfg:     NewConfig(),
	}

	if err := j.ValidateCron(); err != nil {
		t.Fatal(err)
	}

	jr := j.execCommand("test")
	assert.Equal(t, jr.Status, 0)
}

func TestInvalidCron(t *testing.T) {
	j := &JobSpec{
		Cron:    "INVALID",
		Name:    "test",
		Command: []string{"echo", "bar"},
		cfg:     NewConfig(),
	}

	assert.Error(t, j.ValidateCron())

	j = &JobSpec{
		Cron:    "@1minutes",
		Name:    "test",
		Command: []string{"echo", "bar"},
		cfg:     NewConfig(),
	}

	assert.Error(t, j.ValidateCron())
}

func TestJobWithEnvVars(t *testing.T) {
	jobSpec := []byte(`
cron: "* * * * *"
command: env
env: 
  foo: bar
  coffee: bar
`)

	j := JobSpec{}
	err := yaml.Unmarshal(jobSpec, &j)
	if err != nil {
		log.Fatal(err)
	}

	j.ValidateCron()

	_, ok := j.Env["foo"]
	if !ok {
		t.Fatal("should contain foo")
	}

	jr := j.execCommand("test")

	assert.Contains(t, jr.Log, "foo=bar")
}

func TestJobRunNoCommand(t *testing.T) {
	j := &JobSpec{
		Cron: "* * * * *",
		Name: "test",
		cfg:  NewConfig(),
	}

	jr := j.execCommand("test")
	assert.NotEqual(t, jr.Status, 0)
}

func TestJobNonZero(t *testing.T) {
	j := &JobSpec{
		Cron: "* * * * *",
		Name: "test",
		Command: []string{
			"la", "--moo",
		},
		cfg: NewConfig(),
	}

	jr := j.execCommand("test")
	assert.NotEqual(t, jr.Status, 0)
}

func TestJobRunInvalidSchedule(t *testing.T) {
	s := Schedule{}
	j := &JobSpec{
		Cron:    "MooIAmACow",
		Name:    "Bertha",
		Command: []string{"ls"},
		cfg:     NewConfig(),
	}
	s.Jobs = map[string]*JobSpec{}
	s.Jobs["Bertha"] = j

	assert.Error(t, s.Validate())
	// fix cron but add invalid ref
	s.Jobs["Bertha"].Cron = "* * * * *"
	s.Jobs["Bertha"].OnSuccess.TriggerJob = []string{"IDontExist"}

	assert.Error(t, s.Validate())
}
