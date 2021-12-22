package cheek

import (
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
	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	for {
		cursor--
		_, err := fileHandle.Seek(cursor, io.SeekEnd)
		if err != nil {
			return []string{}, err
		}

		char := make([]byte, 1)
		_, err = fileHandle.Read(char)
		if err != nil {
			return []string{}, err
		}

		// nts: char 10 is newline, char 13 is carriage return
		if cursor != -1 && (char[0] == 10 || char[0] == 13) {
			// break
			lines = append(lines, line)
			if nLines > 0 && len(lines) == nLines {
				break
			}
			line = ""

		}

		line = fmt.Sprintf("%s%s", string(char), line)

		if cursor == -filesize { // at beginning of file
			lines = append(lines, line)
			break
		}
	}

	return lines, err
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
