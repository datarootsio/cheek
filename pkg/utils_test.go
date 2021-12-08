package jdi

import (
	"testing"
)

func TestLastLineReader(t *testing.T) {
	l, err := readLastLines("../testdata/test.jsonl", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 2 {
		t.Fatal("incorrect number of lines")
	}
	if l[0] != "{\"a\":4}\n" {
		t.Fatal("incorrect line in first place")
	}

	// go over number of lines in file
	l, err = readLastLines("../testdata/test.jsonl", 20)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) != 4 {
		t.Fatal("incorrect number of lines")
	}
	if l[0] != "{\"a\":4}\n" {
		t.Fatal("incorrect line in first place")
	}

}
