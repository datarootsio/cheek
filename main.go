package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/adhocore/gronx"
	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog/log"
)

type JobsSpec struct {
	Jobs []JobSpec `yaml:"jobs"`
}

type JobSpec struct {
	Name    string `yaml:"name"`
	Cron    string `yaml:"cron"`
	Command string `yaml:"command"`
	startAt time.Time
}

func (j *JobSpec) Do() {
	gronx := gronx.New()

	if gronx.IsValid(j.Cron) != true {
		log.Fatal().Msg("cron string not valid")
	}

	ticker := time.NewTicker(time.Second)

	go func() {
		for range ticker.C {
			due, _ := gronx.IsDue(j.Cron)

			if due {
				out, err := exec.Command(j.Command).Output()
				spew.Dump(out, err)
			}

		}
	}()
}

func getGron(cronExpr string) (gronx.Gronx, error) {
	gron := gronx.New()
	if !gron.IsValid(cronExpr) {
		return gronx.Gronx{}, errors.New("Invalid CRON expression")
	}

	return gron, nil

}

func readSpecs(fn string) (JobsSpec, error) {
	yfile, err := ioutil.ReadFile(fn)

	if err != nil {
		log.Error().Err(err)
		return JobsSpec{}, err
	}

	specs := JobsSpec{}

	err2 := yaml.Unmarshal(yfile, &specs)

	if err2 != nil {

		log.Error().Err(err)
		return JobsSpec{}, err
	}

	return specs, nil

}

func bgTask() {
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		fmt.Println("Tock")
	}
}

func main() {
	fmt.Println("Go Tickers Tutorial")

	go bgTask()

	// This print statement will be executed before
	// the first `tock` prints in the console
	fmt.Println("The rest of my application can continue")

	// here we use an empty select{} in order to keep
	// our main function alive indefinitely as it would
	// complete before our backgroundTask has a chance
	// to execute if we didn't.
	select {}

}
