package cmd

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestConfig(t *testing.T) {
	initConfig()

	vs := []bool{
		viper.IsSet("port"), viper.IsSet("suppressLogs"), viper.IsSet("logLevel"), viper.IsSet("pretty"), viper.IsSet("homedir"),
	}

	for i, v := range vs {
		if !v {
			t.Fatalf("a default viper value has not been set for var with index %v", i)
		}
	}

}

func TestEnvVar(t *testing.T) {
	// check if this works how I assume it works
	initConfig()

	if !viper.GetBool("pretty") {
		t.Fatalf("default value not correct")
	}

	os.Setenv("CHEEK_PRETTY", "false")
	viper.Reset()
	initConfig()
	if viper.GetBool("pretty") {
		t.Fatalf("env var not picked up")
	}

}
