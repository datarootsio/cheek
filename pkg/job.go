package cheek

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/rs/zerolog/log"
)

// JobSpec holds specifications and metadata of a job.
type JobSpec struct {
	Cron           string      `yaml:"cron,omitempty" json:"cron,omitempty"`
	Command        stringArray `yaml:"command" json:"command"`
	Triggers       []string    `yaml:"triggers,omitempty" json:"triggers,omitempty"`
	Name           string      `json:"name"`
	Retries        int         `yaml:"retries,omitempty" json:"retries,omitempty"`
	globalSchedule *Schedule
	runs           []JobRun
}

// JobRun holds information about a job execution.
type JobRun struct {
	Status      int       `json:"status"`
	Log         string    `json:"log"`
	Name        string    `json:"name"`
	TriggeredAt time.Time `json:"triggered_at"`
	TriggeredBy string    `json:"triggered_by"`
	Triggered   []string  `json:"triggered,omitempty"`
}

func (j *JobSpec) loadRuns() {
	const nRuns int = 30
	logFn := path.Join(CheekPath(), fmt.Sprintf("%s.job.jsonl", j.Name))
	jrs, err := readLastJobRuns(logFn, nRuns)
	if err != nil {
		log.Warn().Str("job", j.Name).Err(err).Msgf("could not load job logs from '%s'", logFn)
	}
	j.runs = jrs
}

func (j *JobRun) logToDisk() {
	logFn := path.Join(CheekPath(), fmt.Sprintf("%s.job.jsonl", j.Name))
	f, err := os.OpenFile(logFn,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Warn().Str("job", j.Name).Err(err).Msgf("Can't open job log '%v' for writing", logFn)
		return
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(j); err != nil {
		log.Warn().Str("job", j.Name).Err(err).Msg("Couldn't save job log to disk.")
	}
}

func (j *JobSpec) execCommandWithRetry(trigger string, suppressLogs bool) {
	tries := 0
	var jr JobRun
	const timeOut = 5 * time.Second

	for tries < j.Retries+1 {

		switch {
		case tries == 0:
			jr = j.execCommand(trigger, suppressLogs)
		default:
			jr = j.execCommand(fmt.Sprintf("%s[retry=%v]", trigger, tries), suppressLogs)
		}

		if jr.Status == 0 {
			break
		}
		log.Debug().Str("job", j.Name).Msgf("Job exited unsuccessfully, launching retry after %v timeout.", timeOut)
		tries++
		time.Sleep(timeOut)

	}
}

func (j *JobSpec) execCommand(trigger string, suppressLogs bool) JobRun {
	log.Info().Str("job", j.Name).Str("trigger", trigger).Msgf("Job triggered")
	// init status to non-zero until execution says otherwise
	jr := JobRun{Name: j.Name, TriggeredAt: time.Now(), TriggeredBy: trigger, Status: -1}

	var cmd *exec.Cmd
	switch len(j.Command) {
	case 0:
		err := errors.New("no command specified")
		jr.Log = fmt.Sprintf("Job unable to start: %v", err.Error())
		log.Warn().Str("job", j.Name).Err(err).Msgf("Job unable to start")
		if !suppressLogs {
			fmt.Println(err.Error())
		}
		jr.logToDisk()
		return jr
	case 1:
		cmd = exec.Command(j.Command[0])
	default:
		cmd = exec.Command(j.Command[0], j.Command[1:]...)
	}

	outPipe, _ := cmd.StdoutPipe()
	errPipe, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		jr.Log = fmt.Sprintf("Job unable to start: %v", err.Error())
		if !suppressLogs {
			fmt.Println(err.Error())
		}
		log.Warn().Str("job", j.Name).Err(err).Msgf("Job unable to start")
		jr.logToDisk()
		return jr
	}

	merged := io.MultiReader(outPipe, errPipe)
	reader := bufio.NewReader(merged)
	line, err := reader.ReadString('\n')

	for err == nil {
		if !suppressLogs {
			fmt.Print(line)
		}
		jr.Log += line
		line, err = reader.ReadString('\n')
	}

	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			jr.Status = exitError.ExitCode()
			log.Warn().Str("job", j.Name).Msgf("Exit code %v", exitError.ExitCode())
		}
		return jr
	}

	jr.Status = 0
	// trigger jobs that should run on successful completion
	for _, tn := range j.Triggers {
		tj := j.globalSchedule.Jobs[tn]
		go func(jobName string) {
			tj.execCommandWithRetry(fmt.Sprintf("job[%s]", jobName), suppressLogs)
		}(tn)
	}

	jr.logToDisk()

	return jr
}
