package logger

import (
    "strings"
    "github.com/echocat/caretakerd/logger/level"
    "io"
)

type Receiver struct {
    logger *Logger
    level  level.Level
}

func (i *Logger) ReceiverFor(l level.Level) (*Receiver) {
    return &Receiver{
        logger: i,
        level: l,
    }
}

func (i *Logger) Stdout() (*Receiver) {
    return i.ReceiverFor(i.config.StdoutLevel)
}

func (i *Logger) Stderr() (*Receiver) {
    return i.ReceiverFor(i.config.StderrLevel)
}

func (i *Logger) Stdin() io.Reader {
    return nil
    /*if i.output != nil {
        return nil
    } else {
        return os.Stdin
    }*/
}

func (i Receiver) Write(p []byte) (n int, err error) {
    what := string(p)
    lines := strings.Split(what, "\n")
    numberOfLines := len(lines)
    for j, line := range lines {
        if j < (numberOfLines - 1) || len(line) > 0 {
            i.logger.Log(i.level, line)
        }
    }
    return len(p), nil
}

