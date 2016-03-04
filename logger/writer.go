package logger

import (
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/values"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
)

var writeSynchronizerLock = new(sync.Mutex)
var filenameToWriteSynchronizer = map[values.String]*Writer{}

// Writer represents a log writer.
// A writer of this type is synchronized and could be used from different threads and contexts.
type Writer struct {
	lock       *sync.Mutex
	ownerCount int
	filename   values.String
	writer     *lumberjack.Logger
}

// NewWriter creates a new Write for given filename.
func NewWriter(filename values.String, writer *lumberjack.Logger) *Writer {
	writeSynchronizerLock.Lock()
	defer writeSynchronizerLock.Unlock()

	result := filenameToWriteSynchronizer[filename]
	if result != nil {
		result.lock.Lock()
		result.ownerCount++
		result.lock.Unlock()
	} else {
		result = &Writer{
			lock:       new(sync.Mutex),
			ownerCount: 1,
			filename:   filename,
			writer:     writer,
		}
		filenameToWriteSynchronizer[filename] = result
	}
	return result
}

// Write writes the given content to writer drain.
func (instance *Writer) Write(what []byte, stderr bool) (int, error) {
	instance.lock.Lock()
	defer instance.lock.Unlock()
	writer := instance.writer
	if writer != nil {
		return instance.writeIgnoringProblems(what, writer)
	} else if stderr {
		return os.Stderr.Write(what)
	} else {
		return os.Stdout.Write(what)
	}
}

func (instance *Writer) writeIgnoringProblems(what []byte, to io.Writer) (int, error) {
	var n int
	var err error
	defer func() {
		terr := recover()
		if terr != nil {
			err = errors.New("%v", terr)
		}
	}()
	n, err = to.Write(what)
	return n, err
}

// Close closes this Writer and all of its resources.
// After this a usage is not longer possible.
func (instance *Writer) Close() {
	instance.lock.Lock()
	defer instance.lock.Unlock()
	instance.ownerCount--

	if instance.ownerCount == 0 {
		writeSynchronizerLock.Lock()
		delete(filenameToWriteSynchronizer, instance.filename)
		writeSynchronizerLock.Unlock()
	}
}
