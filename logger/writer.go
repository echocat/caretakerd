package logger

import (
	"sync"
	"io"
	"os"
	. "github.com/echocat/caretakerd/values"
	"github.com/echocat/caretakerd/errors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var writeSynchronizerLock = new(sync.Mutex)
var filenameToWriteSynchronizer = map[String]*Writer{}

type Writer struct {
	lock       *sync.Mutex
	ownerCount int
	filename   String
	writer     *lumberjack.Logger
}

func NewWriteFor(filename String, writer *lumberjack.Logger) *Writer {
	writeSynchronizerLock.Lock()
	defer writeSynchronizerLock.Unlock()

	result := filenameToWriteSynchronizer[filename]
	if result != nil {
		result.lock.Lock()
		result.ownerCount += 1
		result.lock.Unlock()
	} else {
		result = &Writer{
			lock: new(sync.Mutex),
			ownerCount: 1,
			filename: filename,
			writer: writer,
		}
		filenameToWriteSynchronizer[filename] = result
	}
	return result
}

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

func (instance *Writer) Close() {
	instance.lock.Lock()
	defer instance.lock.Unlock()
	instance.ownerCount -= 1

	if instance.ownerCount == 0 {
		writeSynchronizerLock.Lock()
		delete(filenameToWriteSynchronizer, instance.filename)
		writeSynchronizerLock.Unlock()
	}
}
