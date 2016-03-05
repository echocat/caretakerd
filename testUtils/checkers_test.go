package testUtils

import (
	"github.com/echocat/caretakerd/values"
	. "gopkg.in/check.v1"
	"testing"
)

type CheckersTest struct{}

func (s *CheckersTest) TestIsEmpty(c *C) {
	c.Assert("abc", Not(IsEmpty))
	c.Assert("", IsEmpty)
	c.Assert(values.String("abc"), Not(IsEmpty))
	c.Assert(values.String(""), IsEmpty)
	c.Assert([]string{"abc"}, Not(IsEmpty))
	c.Assert([]string{}, IsEmpty)
	c.Assert([]int{1}, Not(IsEmpty))
	c.Assert([]int{}, IsEmpty)
	c.Assert(map[string]int{"abc": 1}, Not(IsEmpty))
	c.Assert(map[string]int{}, IsEmpty)
}

func Test(t *testing.T) {
	TestingT(t)
}

func init() {
	Suite(&CheckersTest{})
}
