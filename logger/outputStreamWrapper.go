package logger

import (
	"io"
	"strings"
)

type outputStreamWrapper struct {
	logger *Logger
	level  Level
}

// NewOutputStreamWrapperFor creates a writer to redirect every output to a logger.
func (i *Logger) NewOutputStreamWrapperFor(l Level) io.Writer {
	return &outputStreamWrapper{
		logger: i,
		level:  l,
	}
}

// Stdout creates a writer to redirect every Stdout output to a logger.
func (i *Logger) Stdout() io.Writer {
	return i.NewOutputStreamWrapperFor(i.config.StdoutLevel)
}

// Stderr creates a writer to redirect every Stderr output to a logger.
func (i *Logger) Stderr() io.Writer {
	return i.NewOutputStreamWrapperFor(i.config.StderrLevel)
}

// Write writes given bytes to logger. It treats every new line as a new log entry.
func (i outputStreamWrapper) Write(p []byte) (n int, err error) {
	what := string(p)
	lines := strings.Split(what, "\n")
	numberOfLines := len(lines)
	for j, line := range lines {
		if j < (numberOfLines-1) || len(line) > 0 {
			i.logger.Log(i.level, line)
		}
	}
	return len(p), nil
}
