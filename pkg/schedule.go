package cheek

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/adhocore/gronx"
	"github.com/rs/zerolog"
)

// Schedule defines specs of a job schedule.
type Schedule struct {
	Jobs       map[string]*JobSpec `yaml:"jobs" json:"jobs"`
	OnSuccess  OnEvent             `yaml:"on_success,omitempty" json:"on_success,omitempty"`
	OnError    OnEvent             `yaml:"on_error,omitempty" json:"on_error,omitempty"`
	TZLocation string              `yaml:"tz_location,omitempty" json:"tz_location,omitempty"`
	loc        *time.Location
	logs       string
	log        zerolog.Logger
	cfg        Config
}

// Run a Schedule based on its specs.
func (s *Schedule) Run() {
	gronx := gronx.New()
	ticker := time.NewTicker(time.Minute)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			for _, j := range s.Jobs {
				if j.Cron == "" {
					continue
				}
				due, _ := gronx.IsDue(j.Cron, s.now())

				if due {
					go func(j *JobSpec) {
						j.execCommandWithRetry("cron")
					}(j)
				}
			}

		case sig := <-sigs:
			s.log.Info().Msgf("%s signal received, exiting...", sig.String())
			return
		}
	}
}

type stringArray []string

func (a *stringArray) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		*a = strings.Fields(single)
	} else {
		*a = multi
	}
	return nil
}

func readSpecs(fn string) (Schedule, error) {
	yfile, err := ioutil.ReadFile(fn)
	if err != nil {
		return Schedule{}, err
	}

	specs := Schedule{}

	if err = yaml.Unmarshal(yfile, &specs); err != nil {
		return Schedule{}, err
	}

	return specs, nil
}

// Validate Schedule spec and logic.
func (s *Schedule) Validate() error {
	for k, v := range s.Jobs {
		// check if trigger references exist
		triggerJobs := append(v.OnSuccess.TriggerJob, v.OnError.TriggerJob...)
		for _, t := range triggerJobs {
			if _, ok := s.Jobs[t]; !ok {
				return fmt.Errorf("cannot find spec of job '%s' that is referenced in job '%s'", t, k)
			}
		}
		// set some metadata & refs for each job
		// for easier retrievability
		v.Name = k
		v.globalSchedule = s
		v.log = s.log
		v.cfg = s.cfg

		// validate cron string
		if err := v.ValidateCron(); err != nil {
			return err
		}

	}
	// validate tz location
	if s.TZLocation == "" {
		s.TZLocation = "Local"
	}

	loc, err := time.LoadLocation(s.TZLocation)
	if err != nil {
		return err
	}
	s.loc = loc

	return nil
}

func (s *Schedule) now() time.Time {
	return time.Now().In(s.loc)
}

func loadSchedule(log zerolog.Logger, cfg Config, fn string) (Schedule, error) {
	s, err := readSpecs(fn)
	if err != nil {
		return Schedule{}, err
	}
	s.log = log
	s.cfg = cfg

	// run validations
	if err := s.Validate(); err != nil {
		return Schedule{}, err
	}

	return s, nil
}

func server(s *Schedule) {
	var httpAddr string = fmt.Sprintf(":%s", s.cfg.Port)
	type Healthz struct {
		Jobs   int    `json:"jobs"`
		Status string `json:"status"`
	}

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		status := Healthz{Jobs: len(s.Jobs), Status: "ok"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	s.log.Info().Msgf("Starting HTTP server on %v", httpAddr)
	s.log.Fatal().Err(http.ListenAndServe(httpAddr, nil))
}

// RunSchedule is the main entry entrypoint of cheek.
func RunSchedule(log zerolog.Logger, cfg Config, scheduleFn string) error {
	s, err := loadSchedule(log, cfg, scheduleFn)
	if err != nil {
		return err
	}
	numberJobs := len(s.Jobs)
	i := 1
	for _, job := range s.Jobs {
		s.log.Info().Msgf("Initializing (%v/%v) job: %s", i, numberJobs, job.Name)
		i++
	}
	go server(&s)
	s.Run()
	return nil
}
