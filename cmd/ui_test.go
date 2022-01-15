package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUICmd(t *testing.T) {
	rootCmd.SetArgs([]string{"ui"}) // gonna fail because server not up
	err := rootCmd.Execute()
	assert.Error(t, err)
}
