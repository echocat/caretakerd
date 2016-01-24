package logger

import (
    . "github.com/echocat/caretakerd/values"
    "sync"
    "io"
    "os"
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

func (this *Writer) Write(what []byte, stderr bool) (int, error) {
    this.lock.Lock()
    defer this.lock.Unlock()
    writer := this.writer
    if writer != nil {
        return this.writeIgnoringProblems(what, writer)
    } else if stderr {
        return os.Stderr.Write(what)
    } else {
        return os.Stdout.Write(what)
    }
}

func (this *Writer) writeIgnoringProblems(what []byte, to io.Writer) (int, error) {
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

func (this *Writer) Close() {
    this.lock.Lock()
    defer this.lock.Unlock()
    this.ownerCount -= 1

    if this.ownerCount == 0 {
        writeSynchronizerLock.Lock()
        delete(filenameToWriteSynchronizer, this.filename)
        writeSynchronizerLock.Unlock()
    }
}
