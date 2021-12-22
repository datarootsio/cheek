package cheek

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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

	// read everything
	l, err = readLastLines("../testdata/test.jsonl", -1)
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

func TestHardWrap(t *testing.T) {
	test := "12345678"
	assert.Equal(t, hardWrap(test, 5), "12345\n678")
	assert.Equal(t, hardWrap(test, 2), "12\n34\n56\n78")
}

func TestCheekPath(t *testing.T) {
	assert.True(t, strings.Contains(CheekPath(), ".cheek"))

	const dirName = "moo_i_am_sheep"

	viper.Set("homedir", dirName)
	assert.True(t, strings.Contains(CheekPath(), "sheep"))

	// cleanup
	os.RemoveAll(dirName)
}
