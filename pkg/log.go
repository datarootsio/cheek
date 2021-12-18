package cheek

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func PrettyStdout() io.Writer {
	return zerolog.ConsoleWriter{Out: os.Stdout}
}

func coreJsonLogger() io.Writer {
	const logFile string = "core.cheek.jsonl"
	logFn := path.Join(cheekPath(), logFile)

	f, err := os.OpenFile(logFn,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Printf("Can't open log file '%s' for writing.", logFile)
		os.Exit(1)
	}
	return f
}

// Configures the package's global logger, also allows to pass in custom writers for
// testing purposes. If no additional loggers specified logs will only be written to the core JSON log.
func ConfigLogger(logLevel string, extraWriters ...io.Writer) {
	var multi zerolog.LevelWriter

	var loggers []io.Writer
	loggers = append(loggers, coreJsonLogger())
	loggers = append(loggers, extraWriters...)

	multi = zerolog.MultiLevelWriter(loggers...)
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		fmt.Printf("Exiting, cannot initialize logger with level '%s'\n", logLevel)
		os.Exit(1)
	}
	log.Logger = zerolog.New(multi).With().Timestamp().Logger().Level(level)
}
