package cheek

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"sync"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const coreLogFile string = "core.cheek.jsonl"

type Config struct {
	Pretty       bool   `yaml:"pretty"`
	SuppressLogs bool   `yaml:"suppressLogs"`
	LogLevel     string `yaml:"logLevel"`
	HomeDir      string `yaml:"homedir"`
	Port         string `yaml:"port"`
	DBPath       string `yaml:"dbpath"`
}

func NewConfig() Config {
	return Config{
		Pretty:       true,
		SuppressLogs: false,
		LogLevel:     "info",
		HomeDir:      CheekPath(),
		Port:         "8081",
		DBPath:       path.Join(CheekPath(), "cheek.sqlite3"),
	}
}

func CheekPath() string {
	var p string
	switch viper.IsSet("homedir") {
	case true:
		p = viper.GetString("homedir")
	default:
		usr, _ := user.Current()
		dir := usr.HomeDir
		p = path.Join(dir, ".cheek")
	}

	_ = os.MkdirAll(p, os.ModePerm)

	return p
}

func readLastJobRuns(log zerolog.Logger, filepath string, nRuns int) ([]JobRun, error) {
	lines, err := readLastLines(filepath, nRuns)
	if err != nil {
		return []JobRun{}, nil
	}

	var jrs []JobRun
	for _, line := range lines {
		jr := JobRun{}
		err = json.Unmarshal([]byte(line), &jr)
		if err != nil {
			log.Debug().Str("logfile", filepath).Err(err).Msgf("can't decode log line: %s", line)
			// try to still fetch other log entries by skipping this log line
			continue
		}
		jrs = append(jrs, jr)
	}

	return jrs, nil
}

func readLastLines(filepath string, nLines int) ([]string, error) {
	fileHandle, err := os.Open(filepath)
	if err != nil {
		return []string{}, err
	}
	defer fileHandle.Close()

	var lines []string
	reader := bufio.NewReader(fileHandle)

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return []string{}, err
		}

		lines = append([]string{s}, lines...)

		if nLines > 0 && len(lines) > nLines {
			lines = lines[:nLines]
		}
	}

	return lines, nil
}

func hardWrap(in string, width int) string {
	if width < 1 {
		return in
	}

	wrapped := ""

	var i int
	for i = 0; len(in[i:]) > width; i += width {
		wrapped += in[i:i+width] + "\n"
	}
	wrapped += in[i:]

	return wrapped
}

type tsBuffer struct {
	b bytes.Buffer
	m sync.Mutex
}

func (b *tsBuffer) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Read(p)
}

func (b *tsBuffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}

func (b *tsBuffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.String()
}

func (b *tsBuffer) Reset() {
	b.m.Lock()
	defer b.m.Unlock()
	b.b.Reset()
}

func PrettyStdout() io.Writer {
	return zerolog.ConsoleWriter{Out: os.Stdout}
}

func CoreJsonLogger() io.Writer {
	logFn := path.Join(CheekPath(), coreLogFile)

	f, err := os.OpenFile(logFn,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Printf("Can't open log file '%s' for writing.", coreLogFile)
		os.Exit(1)
	}
	return f
}

// Configures the package's global logger, also allows to pass in custom writers for
// testing purposes.
func NewLogger(logLevel string, extraWriters ...io.Writer) zerolog.Logger {
	var multi zerolog.LevelWriter

	var loggers []io.Writer
	loggers = append(loggers, CoreJsonLogger())
	loggers = append(loggers, extraWriters...)

	multi = zerolog.MultiLevelWriter(loggers...)
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		fmt.Printf("Exiting, cannot initialize logger with level '%s'\n", logLevel)
		os.Exit(1)
	}
	return zerolog.New(multi).With().Timestamp().Logger().Level(level)
}
