package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/adhocore/gronx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Schedule struct {
	Jobs map[string]*JobSpec `yaml:"jobs"`
}

func (s *Schedule) Run() {
	gronx := gronx.New()
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		for _, j := range s.Jobs {
			if j.Cron == "" {
				continue
			}
			due, _ := gronx.IsDue(j.Cron)
			// spew.Dump(999, due, j.Cron, j.name, j)

			if due {
				go func(j *JobSpec) {
					j.ExecCommand("cron")
				}(j)
			}
		}
	}

}

type JobSpec struct {
	Cron           string   `yaml:"cron"`
	Command        string   `yaml:"command"`
	Triggers       []string `yaml:"triggers"`
	name           string
	globalSchedule *Schedule
	runs           []time.Time
	statuses       []int
	logTails       []string
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
		v.name = k
		v.globalSchedule = &s
	}

	return s, nil
}

func (j *JobSpec) ExecCommand(trigger string) {
	log.Info().Str("job", j.name).Str("trigger", trigger).Msgf("Job triggered")
	// spew.Dump(j)
	// os.Exit(1)
	// register new run
	j.runs = append(j.runs, time.Now())
	j.logTails = append(j.logTails, "")

	cmdString := parseCommand(j.Command)
	var args []string
	if len(cmdString) > 1 {
		args = cmdString[1:]
	}
	cmd := exec.Command(cmdString[0], args...)

	pipe, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		// handle error
	}
	reader := bufio.NewReader(pipe)
	line, err := reader.ReadString('\n')
	for err == nil {
		// output to stdout
		fmt.Print(line)
		// output to our logger
		j.logTails[len(j.logTails)-1] = j.logTails[len(j.logTails)-1] + line + "\n"
		line, err = reader.ReadString('\n')
	}

	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			j.statuses = append(j.statuses, exitError.ExitCode())
		}
		return
	}

	j.statuses = append(j.statuses, 0)
	// trigger jobs that should run on succesful completion
	for _, tn := range j.Triggers {
		tj := j.globalSchedule.Jobs[tn]
		go func(jobName string) {
			tj.ExecCommand(fmt.Sprintf("job[%s]", jobName))
		}(tn)

	}

}

func parseCommand(command string) []string {
	return strings.Fields(command)
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 1 {
		panic("Please pass a schedule file as first argument.")
	}

	js, err := LoadSchedule(argsWithoutProg[0])
	if err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}
	numberJobs := len(js.Jobs)
	i := 0
	for _, job := range js.Jobs {
		log.Info().Msgf("Initializing (%v/%v) job: %s", i, numberJobs, job.name)
		i = i + 1
	}

	js.Run()

}
