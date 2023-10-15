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
	assert.Nil(t, err)
	assert.Equal(t, 2, len(l), "incorrect number of lines")
	assert.Equal(t, "{\"a\":4}\n", l[0], "incorrect line in first place")

	// go over number of lines in file
	l, err = readLastLines("../testdata/test.jsonl", 20)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(l), "incorrect number of lines")
	assert.Equal(t, "{\"a\":4}\n", l[0], "incorrect line in first place")

	// read everything
	l, err = readLastLines("../testdata/test.jsonl", -1)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(l), "incorrect number of lines")
	assert.Equal(t, "{\"a\":4}\n", l[0], "incorrect line in first place")
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
