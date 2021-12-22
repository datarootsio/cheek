package cheek

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/rs/zerolog"
)

// Configures the package's global logger, also allows to pass in custom writers for
// testing purposes.
func NewLogger(prettyLog bool, logLevel string, extraWriters ...io.Writer) zerolog.Logger {
	var multi zerolog.LevelWriter

	const logFile string = "core.cheek.jsonl"
	logFn := path.Join(CheekPath(), logFile)

	f, err := os.OpenFile(logFn,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Printf("Can't open log file '%s' for writing.", logFile)
		os.Exit(1)
	}

	var loggers []io.Writer
	loggers = append(loggers, f)
	loggers = append(loggers, extraWriters...)

	if prettyLog {
		loggers = append(loggers, zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		loggers = append(loggers, os.Stdout)
	}

	multi = zerolog.MultiLevelWriter(loggers...)
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		fmt.Printf("Exiting, cannot initialize logger with level '%s'\n", logLevel)
		os.Exit(1)
	}
	return zerolog.New(multi).With().Timestamp().Logger().Level(level)
}
