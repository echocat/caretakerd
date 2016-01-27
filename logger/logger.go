package logger

import (
	"github.com/echocat/caretakerd/panics"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
	"fmt"
	"sync"
	"time"
	"runtime"
	usync "github.com/echocat/caretakerd/sync"
)

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

func NewLogger(conf Config, name string, syncGroup *usync.SyncGroup) (*Logger, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}
	filename := conf.Filename.String()
	var output *lumberjack.Logger;
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
		config: conf,
		name: name,
		open: true,
		syncGroup: syncGroup,
		lock: new(sync.Mutex),
		output: output,
		created: time.Now(),
		writeSynchronizer: NewWriteFor(conf.Filename, output),
	}
	runtime.SetFinalizer(result, finalize)
	return result, nil
}

func finalize(what *Logger) {
	what.Close()
}

func (i *Logger) Log(level Level, pattern interface{}, args ...interface{}) {
	i.LogCustom(1, nil, level, pattern, args...)
}

func (i *Logger) LogProblem(problem interface{}, level Level, pattern interface{}, args ...interface{}) {
	i.LogCustom(1, problem, level, pattern, args...)
}

func (i *Logger) LogCustom(framesToSkip int, problem interface{}, level Level, pattern interface{}, args ...interface{}) {
	if level >= i.config.Level {
		now := time.Now()
		message := FormatMessage(pattern, args...)
		entry := i.EntryFor(framesToSkip + 1, problem, level, now, message)
		toLog, err := entry.Format(i.config.Pattern, framesToSkip + 1)
		if err != nil {
			panics.New("Could not format log entry with given pattern '%v'. Got: %v", i.config.Pattern, err)
		}
		i.write(level, []byte(toLog))
	}
}

func FormatMessage(pattern interface{}, args ...interface{}) string {
	patternAsString := fmt.Sprintf("%v", pattern)
	if len(args) > 0 {
		//noinspection GoPlaceholderCount
		return fmt.Sprintf(patternAsString, args...)
	} else {
		return patternAsString
	}
}

func (i *Logger) write(level Level, message []byte) {
	i.lock.Lock()
	defer i.unlocker()
	if ! i.IsOpen() {
		panics.New("The logger is not open.").Throw()
	}
	i.writeSynchronizer.Write(message, level.IsIndicatingProblem())
}

func (i Logger) IsOpen() bool {
	return i.open
}

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

func (i Logger) Uptime() time.Duration {
	return time.Since(i.created)
}
