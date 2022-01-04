package cheek

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/adhocore/gronx"
	"github.com/rs/zerolog"
)

// OnEvent contains specs on what needs to happen after a job event.
type OnEvent struct {
	TriggerJob    []string `yaml:"trigger_job,omitempty" json:"trigger_job,omitempty"`
	NotifyWebhook []string `yaml:"notify_webhook,omitempty" json:"notify_webhook,omitempty"`
}

// JobSpec holds specifications and metadata of a job.
type JobSpec struct {
	Cron    string      `yaml:"cron,omitempty" json:"cron,omitempty"`
	Command stringArray `yaml:"command" json:"command"`

	OnSuccess OnEvent `yaml:"on_success,omitempty" json:"on_success,omitempty"`
	OnError   OnEvent `yaml:"on_error,omitempty" json:"on_error,omitempty"`

	Name           string `json:"name"`
	Retries        int    `yaml:"retries,omitempty" json:"retries,omitempty"`
	Env            map[string]string
	globalSchedule *Schedule
	runs           []JobRun

	log zerolog.Logger
	cfg Config
}

// JobRun holds information about a job execution.
type JobRun struct {
	Status      int `json:"status"`
	logBuf      bytes.Buffer
	Log         string    `json:"log"`
	Name        string    `json:"name"`
	TriggeredAt time.Time `json:"triggered_at"`
	TriggeredBy string    `json:"triggered_by"`
	Triggered   []string  `json:"triggered,omitempty"`
	jobRef      *JobSpec
}

func (jr *JobRun) Close() {
	jr.Log = jr.logBuf.String()
}

func (j *JobRun) logToDisk() {
	logFn := path.Join(CheekPath(), fmt.Sprintf("%s.job.jsonl", j.Name))
	f, err := os.OpenFile(logFn,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		j.jobRef.log.Warn().Str("job", j.Name).Err(err).Msgf("Can't open job log '%v' for writing", logFn)
		return
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(j); err != nil {
		j.jobRef.log.Warn().Str("job", j.Name).Err(err).Msg("Couldn't save job log to disk.")
	}
}

func (j *JobSpec) execCommandWithRetry(trigger string) {
	tries := 0
	var jr JobRun
	const timeOut = 5 * time.Second

	for tries < j.Retries+1 {

		switch {
		case tries == 0:
			jr = j.execCommand(trigger)
		default:
			jr = j.execCommand(fmt.Sprintf("%s[retry=%v]", trigger, tries))
		}

		if jr.Status == 0 {
			break
		}
		j.log.Debug().Str("job", j.Name).Msgf("Job exited unsuccessfully, launching retry after %v timeout.", timeOut)
		tries++
		time.Sleep(timeOut)

	}
}

func (j *JobSpec) execCommand(trigger string) JobRun {
	j.log.Info().Str("job", j.Name).Str("trigger", trigger).Msgf("Job triggered")
	// init status to non-zero until execution says otherwise
	jr := JobRun{Name: j.Name, TriggeredAt: time.Now(), TriggeredBy: trigger, Status: -1, jobRef: j}

	suppressLogs := j.cfg.SuppressLogs

	defer j.OnEvent(&jr, suppressLogs)
	defer jr.logToDisk()
	defer jr.Close()

	if j.cfg.Telemetry {
		go func() {
			_, err := ET{}.PhoneHome(j.cfg.PhoneHomeUrl)
			if err != nil {
				j.log.Debug().Str("telemetry", "ET").Err(err).Msg("cannot phone home")
			}
		}()
	}

	var cmd *exec.Cmd
	switch len(j.Command) {
	case 0:
		err := errors.New("no command specified")
		jr.Log = fmt.Sprintf("Job unable to start: %v", err.Error())
		j.log.Warn().Str("job", j.Name).Err(err).Msg("Job unable to start")
		if !suppressLogs {
			fmt.Println(err.Error())
		}
		return jr
	case 1:
		cmd = exec.Command(j.Command[0])
	default:
		cmd = exec.Command(j.Command[0], j.Command[1:]...)
	}

	// add env vars
	cmd.Env = os.Environ()
	for k, v := range j.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	var w io.Writer
	switch j.cfg.SuppressLogs {
	case true:
		w = &jr.logBuf
	default:
		w = io.MultiWriter(os.Stdout, &jr.logBuf)
	}

	// merge stdout and stderr to same writer
	cmd.Stdout = w
	cmd.Stderr = w

	err := cmd.Start()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Job unable to start: %v", err.Error())))
		if !suppressLogs {
			fmt.Println(err.Error())
		}
		j.log.Warn().Str("job", j.Name).Err(err).Msg("Job unable to start")

		return jr
	}

	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			jr.Status = exitError.ExitCode()
			j.log.Warn().Str("job", j.Name).Msgf("Exit code %v", exitError.ExitCode())
		}

		return jr
	}

	jr.Status = 0

	return jr
}

func (j *JobSpec) loadRuns() {
	const nRuns int = 30
	logFn := path.Join(CheekPath(), fmt.Sprintf("%s.job.jsonl", j.Name))
	jrs, err := readLastJobRuns(j.log, logFn, nRuns)
	if err != nil {
		j.log.Warn().Str("job", j.Name).Err(err).Msgf("could not load job logs from '%s'", logFn)
	}
	j.runs = jrs
}

func (j *JobSpec) ValidateCron() error {
	if j.Cron != "" {
		gronx := gronx.New()
		if !gronx.IsValid(j.Cron) {
			return fmt.Errorf("cron string for job '%s' not valid", j.Name)
		}
	}
	return nil
}

func (j *JobSpec) OnEvent(jr *JobRun, suppressLogs bool) {
	var jobsToTrigger []string
	var webhooksToCall []string

	switch jr.Status == 0 {
	case true: // after success
		jobsToTrigger = j.OnSuccess.TriggerJob
		webhooksToCall = j.OnSuccess.NotifyWebhook
	case false: // after error
		jobsToTrigger = j.OnError.TriggerJob
		webhooksToCall = j.OnError.NotifyWebhook
	}

	for _, tn := range jobsToTrigger {
		tj := j.globalSchedule.Jobs[tn]
		j.log.Debug().Str("job", j.Name).Str("on_event", "job_trigger").Msg("triggered by parent job")
		go func() {
			tj.execCommandWithRetry(fmt.Sprintf("job[%s]", j.Name))
		}()
	}

	// trigger webhooks
	for _, wu := range webhooksToCall {
		j.log.Debug().Str("job", j.Name).Str("on_event", "webhook_call").Msg("triggered by parent job")
		go func(webhookURL string) {
			resp_body, err := JobRunWebhookCall(jr, webhookURL)
			if err != nil {
				j.log.Warn().Str("job", j.Name).Str("on_event", "webhook").Err(err).Msg("webhook notify failed")
			}
			j.log.Debug().Str("job", jr.Name).Str("webhook_call", "response").Str("webhook_url", webhookURL).Msg(string(resp_body))
		}(wu)
	}
}
