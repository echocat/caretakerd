package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestEntry_Format(t *testing.T) {
	e := NewEntry(
		0,
		nil,
		"My.Cool.Logger",
		Info,
		time.Unix(1452347970, 0).UTC(),
		"This is a test!",
		time.Duration(10*time.Second),
	)
	f, err := e.Format("%d{YYYY-MM-DD HH:mm:ss}/%r [%-5.5p] [%-8.8c{2}] %m at %C{1}.%M(%F{1}:%L)%n", 0)
	assert.Nil(t, err)
	assert.Equal(t, "2016-01-09 13:59:30/10000 [INFO ] [Cool.Log] This is a test! at logger.TestEntry_Format(entry_test.go:10)\n", f)
}
