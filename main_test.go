package main

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func TestCronParser(t *testing.T) {
	s, err := readSpecs("./testdata/jobs1.yaml")
	spew.Dump(s, err)
	if err != nil {
		t.Error(err)
	}

	s.Jobs[0].Do()
	time.Sleep(time.Minute * 5)
}
