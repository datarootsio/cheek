package cheek

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCheekPath(t *testing.T) {
	assert.True(t, strings.Contains(CheekPath(), ".cheek"))

	const dirName = "moo_i_am_sheep"

	viper.Set("homedir", dirName)
	assert.True(t, strings.Contains(CheekPath(), "sheep"))

	// cleanup
	os.RemoveAll(dirName)
}
