package logger

import (
	"fmt"
	"github.com/echocat/caretakerd/panics"
	usync "github.com/echocat/caretakerd/sync"
	"gopkg.in/natefinch/lumberjack.v2"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Logger represents a logger to log events to different sources like console or files.
type Logger struct {
	config            Config
	name              string
	lock              *sync.Mutex
	open              bool
	syncGroup         *usync.SyncGroup
	output            *lumberjack.Logger
	created           time.Time
	writeSynchronizer *Writer
}

// NewLogger creates a new instance of Logger.
func NewLogger(conf Config, name string, syncGroup *usync.SyncGroup) (*Logger, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}
	filename := conf.Filename.String()
	var output *lumberjack.Logger
	if len(strings.TrimSpace(filename)) > 0 && strings.ToLower(filename) != "console" {
		output = &lumberjack.Logger{
			Filename:   conf.Filename.String(),
			MaxSize:    conf.MaxSizeInMb.Int(),
			MaxBackups: conf.MaxBackups.Int(),
			MaxAge:     conf.MaxAgeInDays.Int(),
		}
		output.Rotate()
	} else {
		output = nil
	}
	result := &Logger{
		config:            conf,
		name:              name,
		open:              true,
		syncGroup:         syncGroup,
		lock:              new(sync.Mutex),
		output:            output,
		created:           time.Now(),
		writeSynchronizer: NewWriter(conf.Filename, output),
	}
	runtime.SetFinalizer(result, finalize)
	return result, nil
}

func finalize(what *Logger) {
	what.Close()
}

// Log logs the given pattern with the given level.
func (i *Logger) Log(level Level, pattern interface{}, args ...interface{}) {
	i.LogAdvanced(1, nil, level, pattern, args...)
}

// LogProblem logs a problem with the given pattern and level.
func (i *Logger) LogProblem(problem interface{}, level Level, pattern interface{}, args ...interface{}) {
	i.LogAdvanced(1, problem, level, pattern, args...)
}

// LogAdvanced logs a problem with the given pattern and level.
func (i *Logger) LogAdvanced(framesToSkip int, problem interface{}, level Level, pattern interface{}, args ...interface{}) {
	if level >= i.config.Level {
		now := time.Now()
		message := formatMessage(pattern, args...)
		entry := i.EntryFor(framesToSkip+1, problem, level, now, message)
		toLog, err := entry.Format(i.config.Pattern, framesToSkip+1)
		if err != nil {
			panics.New("Could not format log entry with given pattern '%v'. Got: %v", i.config.Pattern, err)
		}
		i.write(level, []byte(toLog))
	}
}

func formatMessage(pattern interface{}, args ...interface{}) string {
	patternAsString := fmt.Sprintf("%v", pattern)
	if len(args) > 0 {
		//noinspection GoPlaceholderCount
		return fmt.Sprintf(patternAsString, args...)
	}
	return patternAsString
}

func (i *Logger) write(level Level, message []byte) {
	i.lock.Lock()
	defer i.unlocker()
	if !i.IsOpen() {
		panics.New("The logger is not open.").Throw()
	}
	i.writeSynchronizer.Write(message, level.IsIndicatingProblem())
}

// IsOpen returns true if the current logger is still open and usable.
func (i Logger) IsOpen() bool {
	return i.open
}

// Close will close this logger and all of its resources.
func (i *Logger) Close() {
	i.lock.Lock()
	defer i.unlocker()
	defer func() {
		i.output = nil
		i.open = false
	}()
	if i.output != nil {
		i.output.Close()
	}
	i.writeSynchronizer.Close()
}

func (i *Logger) unlocker() {
	i.lock.Unlock()
}

// Uptime returns the uptime duration of this logger.
func (i Logger) Uptime() time.Duration {
	return time.Since(i.created)
}

// EntryFor creates a new entry for given parameters using the current Logger instance.
func (i *Logger) EntryFor(framesToSkip int, problem interface{}, priority Level, time time.Time, message string) Entry {
	uptime := i.Uptime()
	return NewEntry(framesToSkip+1, problem, i.name, priority, time, message, uptime)
}
