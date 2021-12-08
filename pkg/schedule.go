package jdi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/adhocore/gronx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Schedule struct {
	Jobs map[string]*JobSpec `yaml:"jobs" json:"jobs"`
}

func GetStateFromDisk() (*Schedule, error) {
	const jdiStateFile = "jdi.scheduler.json"
	usr, _ := user.Current()
	dir := usr.HomeDir
	stateFn := path.Join(dir, jdiStateFile)

	s := Schedule{}

	f, err := os.Open(stateFn)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	return &s, nil

}

func (s *Schedule) PersistToDisk() {
	const jdiStateFile = ".jdi.json"
	usr, _ := user.Current()
	dir := usr.HomeDir
	stateFn := path.Join(dir, jdiStateFile)

	f, _ := json.MarshalIndent(s, "", " ")
	_ = ioutil.WriteFile(stateFn, f, 0644)
}

func (s *Schedule) Run() {
	gronx := gronx.New()
	ticker := time.NewTicker(time.Second)

	changeNotify := make(chan bool)
	go func(ch <-chan bool) {
		for _ = range ch {
			s.PersistToDisk()
		}
	}(changeNotify)

	for range ticker.C {
		for _, j := range s.Jobs {
			if j.Cron == "" {
				continue
			}
			due, _ := gronx.IsDue(j.Cron)

			if due {
				go func(j *JobSpec) {
					j.ExecCommand("cron")
					changeNotify <- true
				}(j)
			}
		}
	}

}

type JobSpec struct {
	Cron           string      `yaml:"cron,omitempty" json:"cron,omitempty"`
	Command        StringArray `yaml:"command" json:"command"`
	Triggers       []string    `yaml:"triggers,omitempty" json:"triggers,omitempty"`
	Name           string      `json:"name"`
	globalSchedule *Schedule
}

type JobRun struct {
	Status      int       `json:"status"`
	Log         string    `json:"log"`
	Name        string    `json:"name"`
	TriggeredAt time.Time `json:"triggered_at"`
	TriggeredBy string    `json:"triggered_by"`
	Triggered   []string  `json:"triggered,omitempty"`
}

func JdiPath() string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	p := path.Join(dir, ".jdi")
	_ = os.MkdirAll(p, os.ModePerm)

	return p
}

func (j *JobRun) LogToDisk() {
	logFn := path.Join(JdiPath(), fmt.Sprintf("%s.job.jsonl", j.Name))
	f, err := os.OpenFile(logFn,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Warn().Str("job", j.Name).Err(err).Msgf("Can't open job log '%v' for writing", logFn)
		return
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(j); err != nil {
		log.Warn().Str("job", j.Name).Err(err).Msg("Couldn't save job log to disk.")
	}
}

type StringArray []string

func (a *StringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

func readSpecs(fn string) (Schedule, error) {
	yfile, err := ioutil.ReadFile(fn)

	if err != nil {
		log.Error().Err(err)
		return Schedule{}, err
	}

	specs := Schedule{}

	if err = yaml.Unmarshal(yfile, &specs); err != nil {

		log.Error().Err(err)
		return Schedule{}, err
	}

	return specs, nil

}

func LoadSchedule(fn string) (Schedule, error) {
	s, err := readSpecs(fn)
	if err != nil {
		return Schedule{}, err
	}

	// run validations
	for k, v := range s.Jobs {
		// validate cron string
		if v.Cron != "" {
			gronx := gronx.New()
			if !gronx.IsValid(v.Cron) {
				return Schedule{}, fmt.Errorf("cron string for job '%s' not valid", k)

			}
		}
		// check if trigger references exist
		for _, t := range v.Triggers {
			if _, ok := s.Jobs[t]; !ok {
				return Schedule{}, fmt.Errorf("cannot find spec of job '%s' that is referenced in job '%s'", t, k)
			}

		}
		// set name for easy access
		v.Name = k
		v.globalSchedule = &s
	}

	return s, nil
}

func (j *JobSpec) ExecCommand(trigger string) {
	log.Info().Str("job", j.Name).Str("trigger", trigger).Msgf("Job triggered")
	jr := JobRun{Name: j.Name, TriggeredAt: time.Now(), TriggeredBy: trigger}
	cmd := exec.Command(j.Command[0], j.Command[1:]...)

	outPipe, _ := cmd.StdoutPipe()
	errPipe, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		jr.Status = -1
		jr.Log = fmt.Sprintf("Job unable to start: %v", err.Error())
		fmt.Println(err.Error())
		log.Warn().Str("job", j.Name).Err(err).Msgf("Job unable to start")
		return
	}

	merged := io.MultiReader(outPipe, errPipe)
	reader := bufio.NewReader(merged)
	line, err := reader.ReadString('\n')

	for err == nil {
		// output to stdout
		fmt.Print(line)
		jr.Log += line
		line, err = reader.ReadString('\n')
	}

	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// j.Statuses = append(j.Statuses, exitError.ExitCode())
			jr.Status = exitError.ExitCode()
			log.Warn().Str("job", j.Name).Msgf("Exit code %v", exitError.ExitCode())
		}

		return
	}

	jr.Status = 0
	// trigger jobs that should run on succesful completion
	for _, tn := range j.Triggers {
		tj := j.globalSchedule.Jobs[tn]
		go func(jobName string) {
			tj.ExecCommand(fmt.Sprintf("job[%s]", jobName))
		}(tn)
	}

	jr.LogToDisk()

}

func server(s *Schedule) {
	const httpAddr string = ":8080"
	type Healthz struct {
		Jobs   int    `json:"jobs"`
		Status string `json:"status"`
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		status := Healthz{Jobs: len(s.Jobs), Status: "ok"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)

	})

	http.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	})

	log.Info().Msgf("Starting HTTP server on %v", httpAddr)
	log.Fatal().Err(http.ListenAndServe(":8081", nil))

}

func RunSchedule(fn string, prettyLog bool) {
	// config logger
	// log to st
	const logFile string = "core.jdi.jsonl"
	logFn := path.Join(JdiPath(), logFile)
	f, err := os.OpenFile(logFn,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Warn().Err(err).Msgf("Can't open log file '%s' for writing.", logFile)
		return
	}
	defer f.Close()

	var multi zerolog.LevelWriter

	if prettyLog {
		multi = zerolog.MultiLevelWriter(f, zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		multi = zerolog.MultiLevelWriter(f, os.Stdout)
	}

	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	js, err := LoadSchedule(fn)
	if err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}
	numberJobs := len(js.Jobs)
	i := 0
	for _, job := range js.Jobs {
		log.Info().Msgf("Initializing (%v/%v) job: %s", i, numberJobs, job.Name)
		i = i + 1
	}
	go server(&js)
	js.Run()

}
