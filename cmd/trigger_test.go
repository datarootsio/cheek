package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTriggerCmd(t *testing.T) {
	rootCmd.SetArgs([]string{"trigger", "../testdata/jobs1.yaml", "bar"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}
