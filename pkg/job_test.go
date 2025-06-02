package cheek

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLoadLogs(t *testing.T) {
	db, err := OpenDB("./tmp.sqlite3")
	if err != nil {
		t.Fatal(err)
	}
	cfg := NewConfig()
	cfg.DB = db

	l := NewLogger("debug", nil, os.Stdout, os.Stdout)

	j := &JobSpec{
		Cron:    "* * * * *",
		Name:    "test",
		Command: []string{"echo", "bar"},
		cfg:     cfg,
		log:     l,
	}

	_, err = j.ToYAML(false)
	if err != nil {
		t.Fatal(err)
	}

	_ = j.execCommandWithRetry("test")

	// log loading goes on job name basis
	// let's recreate and see if we can load logs

	j = &JobSpec{
		Name: "test",
		cfg:  cfg,
		log:  l,
	}

	j.loadRunsFromDb(10, false)

	assert.Greater(t, len(j.Runs), 0)
}

func TestJobRun(t *testing.T) {
	j := &JobSpec{
		Cron:    "* * * * *",
		Name:    "test",
		Command: []string{"echo", "bar"},
		cfg:     NewConfig(),
	}

	jobRun := JobRun{}

	// Execute command and get result
	jr := j.execCommand(jobRun, "test")

	// Dereference the pointer and compare the value
	assert.Equal(t, *jr.Status, 0)
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

	jobRun := JobRun{}                  // Create a JobRun instance
	jr := j.execCommand(jobRun, "test") // Pass JobRun instance and "test"
	assert.Equal(t, *jr.Status, 0)
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

	if err := j.ValidateCron(); err != nil {
		t.Fatal(err)
	}

	_, ok := j.Env["foo"]
	if !ok {
		t.Fatal("should contain foo")
	}

	jobRun := JobRun{}                  // Create a JobRun instance
	jr := j.execCommand(jobRun, "test") // Pass JobRun instance and "test"

	jr.flushLogBuffer()

	assert.Contains(t, jr.Log, "foo=bar")
}

func TestStdErrOut(t *testing.T) {
	cfg := NewConfig()
	cfg.SuppressLogs = true

	j := &JobSpec{
		Cron:    "* * * * *",
		Name:    "test",
		Command: []string{"sh", "-c", "echo stdout; echo stderr 1>&2"},
		//  1>&2 sends to stderr
		cfg: cfg,
	}

	jobRun := JobRun{}                  // Create a JobRun instance
	jr := j.execCommand(jobRun, "test") // Pass JobRun instance and "test"
	jr.flushLogBuffer()
	assert.Contains(t, jr.Log, "stdout")
	assert.Contains(t, jr.Log, "stderr")
}

func TestFailingLog(t *testing.T) {
	cfg := NewConfig()
	cfg.SuppressLogs = true

	j := &JobSpec{
		Cron:    "* * * * *",
		Name:    "test",
		Command: []string{"this fails"},
		//  1>&2 sends to stderr
		cfg: cfg,
	}

	jobRun := JobRun{}                  // Create a JobRun instance
	jr := j.execCommand(jobRun, "test") // Pass JobRun instance and "test"
	jr.flushLogBuffer()
	assert.Contains(t, jr.Log, "this fails")
}

func TestJobRunNoCommand(t *testing.T) {
	j := &JobSpec{
		Cron: "* * * * *",
		Name: "test",
		cfg:  NewConfig(),
	}

	jobRun := JobRun{}                  // Create a JobRun instance
	jr := j.execCommand(jobRun, "test") // Pass JobRun instance and "test"
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

	jobRun := JobRun{}                  // Create a JobRun instance
	jr := j.execCommand(jobRun, "test") // Pass JobRun instance and "test"
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

	assert.Error(t, s.initialize())
	// fix cron but add invalid ref
	s.Jobs["Bertha"].Cron = "* * * * *"
	s.Jobs["Bertha"].OnSuccess.TriggerJob = []string{"IDontExist"}

	assert.Error(t, s.initialize())
}

func TestOnEventWebhook(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		// mirror this
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, string(body))
	}))

	defer testServer.Close()

	j := &JobSpec{
		Command: []string{"echo"},
		cfg:     NewConfig(),
		OnSuccess: OnEvent{
			NotifyWebhook: []string{testServer.URL},
		},
	}
	jobRun := JobRun{}                  // Create a JobRun instance
	jr := j.execCommand(jobRun, "test") // Pass JobRun instance and "test"
	j.OnEvent(&jr)
}

func TestStringArray(t *testing.T) {
	type testCase struct {
		yamlString         string
		expectedStatus     int
		expectedLogContent string
	}

	for _, scenario := range []testCase{
		{
			yamlString:     `command: echo foo`,
			expectedStatus: 0, expectedLogContent: "foo",
		},
		{
			yamlString: `command:
- echo
- foo`,
			expectedStatus: 0, expectedLogContent: "foo",
		},
	} {
		j := JobSpec{}
		err := yaml.Unmarshal([]byte(scenario.yamlString), &j)
		if err != nil {
			t.Fatal(err)
		}

		j.cfg = NewConfig()
		jobRun := JobRun{}                  // Create a JobRun instance
		jr := j.execCommand(jobRun, "test") // Pass JobRun instance and "test"

		jr.flushLogBuffer()
		assert.Equal(t, *jr.Status, scenario.expectedStatus)
		assert.Contains(t, jr.Log, scenario.expectedLogContent)
	}
}

func TestStandaloneJobRun(t *testing.T) {
	b := new(tsBuffer)
	log := NewLogger("debug", nil, b, os.Stdout)
	cfg := NewConfig()

	jr, err := RunJob(log, cfg, "../testdata/jobs1.yaml", "bar")
	assert.NoError(t, err)
	assert.Contains(t, b.String(), "\"job\":\"bar\",\"trigger\":\"manual\"")
	assert.Contains(t, jr.Log, "bar_foo")
}

func TestWorkingDir(t *testing.T) {
	b := new(tsBuffer)
	log := NewLogger("debug", nil, b, os.Stdout)
	cfg := NewConfig()

	jr, err := RunJob(log, cfg, "../testdata/readme_example.yaml", "other_workingdir")
	assert.NoError(t, err)
	assert.Contains(t, jr.Log, "/testdata")
}

func TestJobWithBashEval(t *testing.T) {
	b := new(tsBuffer)
	log := NewLogger("debug", nil, b, os.Stdout)
	cfg := NewConfig()

	j := &JobSpec{
		Cron: "* * * * *",
		Name: "test",
		Command: []string{
			"bash", "-c", "MY_VAR=$(date +%Y-%m); echo $MY_VAR $FOO",
		},
		// add random Env to check if it passes through
		Env: map[string]secret{
			"FOO": "BAR",
		},
		cfg: NewConfig(),
	}

	j.log = log
	j.cfg = cfg

	jobRun := JobRun{}                  // Create a JobRun instance
	jr := j.execCommand(jobRun, "test") // Pass JobRun instance and "test"
	jr.flushLogBuffer()

	currentYearMonth := time.Now().Format("2006-01")
	assert.Contains(t, jr.Log, currentYearMonth)
	assert.Contains(t, jr.Log, "BAR")
}

// TestExecCommandStartError simulates a scenario where cmd.Start() fails
func TestExecCommandStartError(t *testing.T) {
	// Create a sample JobRun instance
	jobRun := JobRun{
		LogEntryId:  1,
		Status:      nil,
		logBuf:      bytes.Buffer{},
		Log:         "",
		Name:        "TestJob",
		TriggeredAt: time.Now(),
		TriggeredBy: "manual",
		Triggered:   []string{"manual"},
		Duration:    0,
	}

	// Create a sample JobSpec instance with a command that will fail
	jobSpec := JobSpec{
		Name:    "TestJob",
		Command: stringArray{"nonexistent-command"},
		cfg: Config{
			SuppressLogs: false,
		},
	}

	// Run the execCommand method with the JobSpec and JobRun
	trigger := "manual"
	result := jobSpec.execCommand(jobRun, trigger)

	// Assertions
	assert.NotNil(t, result.Status, "Expected job run status to be set")
	assert.Equal(t, StatusError, *result.Status, "Expected StatusError when cmd.Start fails")
	assert.Contains(t, result.Log, "job unable to start", "Expected log to contain failure message")
	assert.NotEmpty(t, result.Log, "Expected log to contain some content")
	assert.NotNil(t, result.Duration, "Expected duration to be set")
	assert.GreaterOrEqual(t, result.Duration.Milliseconds(), int64(0), "Expected positive duration")
}

func TestExecCommandExitError(t *testing.T) {
	// Create a sample JobRun instance
	jobRun := JobRun{
		LogEntryId:  1,
		Status:      nil,
		logBuf:      bytes.Buffer{},
		Log:         "",
		Name:        "TestJob",
		TriggeredAt: time.Now(),
		TriggeredBy: "manual",
		Triggered:   []string{"manual"},
		Duration:    0,
	}

	// Create a sample JobSpec instance with a command that will fail
	jobSpec := JobSpec{
		Name:    "TestJob",
		Command: stringArray{"false"}, // Use a command that always exits with code 1
		cfg: Config{
			SuppressLogs: false,
		},
	}

	// Run the execCommand method with the JobSpec and JobRun
	trigger := "manual"
	result := jobSpec.execCommand(jobRun, trigger)

	// Assertions
	assert.NotNil(t, result.Status, "Expected job run status to be set")
	assert.Equal(t, 1, *result.Status, "Expected StatusError when cmd.Wait fails with non-zero exit code")
	assert.Contains(t, result.Log, "Exit code:", "Expected log to contain exit status message")
	assert.NotEmpty(t, result.Log, "Expected log to contain some content")
	assert.NotNil(t, result.Duration, "Expected duration to be set")
	assert.GreaterOrEqual(t, result.Duration.Milliseconds(), int64(0), "Expected positive duration")
}

func TestMarshalSecret(t *testing.T) {
	secrets := map[string]secret{"foo": "bar"}

	yamlResult, err := yaml.Marshal(secrets)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(yamlResult), "foo: '***'\n")

	jsonResult, err := json.Marshal(secrets)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(jsonResult), `{"foo":"***"}`)
}
