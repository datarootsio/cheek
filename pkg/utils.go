package cheek

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

const jobNameCoreProcess = "_cheek"

type Config struct {
	Pretty       bool   `yaml:"pretty"`
	SuppressLogs bool   `yaml:"suppressLogs"`
	LogLevel     string `yaml:"logLevel"`
	HomeDir      string `yaml:"homedir"`
	Port         string `yaml:"port"`
	DBPath       string `yaml:"dbpath"`
	DB           *sqlx.DB
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

func (c *Config) Init() error {
	var err error
	c.DB, err = OpenDB(c.DBPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	return nil
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

func PrettyStdout() io.Writer {
	return zerolog.ConsoleWriter{Out: os.Stdout}
}

type DBLogWriter struct {
	db *sqlx.DB
}

func (w DBLogWriter) Write(p []byte) (n int, err error) {
	_, err = w.db.Exec("INSERT INTO log (job, message) VALUES (?, ?)", jobNameCoreProcess, string(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func NewDBLogWriter(db *sqlx.DB) io.Writer {
	return DBLogWriter{db: db}
}

// Configures the package's global logger, also allows to pass in custom writers for
// testing purposes.
func NewLogger(logLevel string, db *sqlx.DB, extraWriters ...io.Writer) zerolog.Logger {
	var multi zerolog.LevelWriter

	var loggers []io.Writer
	if db != nil {
		loggers = append(loggers, NewDBLogWriter(db))
	}
	loggers = append(loggers, extraWriters...)

	multi = zerolog.MultiLevelWriter(loggers...)
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		fmt.Printf("Exiting, cannot initialize logger with level '%s'\n", logLevel)
		os.Exit(1)
	}
	return zerolog.New(multi).With().Timestamp().Logger().Level(level)
}

func getCoreLogsFromDB(db *sqlx.DB, nruns int) ([]JobRun, error) {
	var logs []JobRun
	if err := db.Select(&logs, "SELECT triggered_at, message FROM log WHERE job = ? ORDER BY id DESC LIMIT ?", jobNameCoreProcess, nruns); err != nil {
		return nil, err
	}
	return logs, nil
}
