package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	// wrong fn will stop the scheduler from going live
	rootCmd.SetArgs([]string{"run", "../testdata/not-exists.yaml"})
	err := rootCmd.Execute()
	assert.Contains(t, err.Error(), "no such file or directory")
}
