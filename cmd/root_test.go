package cmd

import (
	"os"
	"testing"

	cheek "github.com/bart6114/cheek/pkg"
	"github.com/spf13/viper"
)

func TestEnvVar(t *testing.T) {
	// check if this works how I assume it works
	initConfig()

	if !viper.GetBool("pretty") {
		t.Fatalf("default value not correct")
	}

	_ = os.Setenv("CHEEK_PRETTY", "false")
	initConfig()
	if viper.GetBool("pretty") {
		t.Fatalf("env var not picked up")
	}
}

func TestUnmarshall(t *testing.T) {
	c := cheek.NewConfig()
	if err := viper.Unmarshal(&c); err != nil {
		t.Fatal(err)
	}
}
