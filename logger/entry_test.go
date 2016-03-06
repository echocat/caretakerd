package logger

import (
	. "gopkg.in/check.v1"
	"time"
)

type EntryTest struct{}

func init() {
	Suite(&EntryTest{})
}

func (s *EntryTest) TestFormat(c *C) {
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
	c.Assert(err, IsNil)
	c.Assert(f, Matches, "2016-01-09 13:59:30/10000 \\[INFO \\] \\[Cool\\.Log\\] This is a test! at .+\\.TestFormat\\(entry_test.go:[0-9]+\\)\n")
}
