package logger

import (
    "testing"
    "time"
    "github.com/echocat/caretakerd/logger/level"
    "github.com/stretchr/testify/assert"
)

func TestEntry_Format(t *testing.T) {
    e := NewEntry(
        0,
        "My.Cool.Logger",
        level.Info,
        time.Unix(1452347970, 0),
        "This is a test!",
        time.Duration(10 * time.Second),
    )
    f, err := e.Format("%d{YYYY-MM-DD HH:mm:ss}/%r [%-5.5p] [%-8.8c{2}] %m at %C{1}.%M(%F{1}:%L)%n")
    assert.Nil(t, err)
    assert.Equal(t, "2016-01-09 14:59:30/10000 [INFO ] [Cool.Log] This is a test! at logger.TestEntry_Format(entry_test.go:11)\n", f)
}
