package caretakerd

import (
	. "github.com/echocat/caretakerd/values"
	. "gopkg.in/check.v1"
	"testing"
)

type ConfigEnvironmentTest struct{}

func (s *ConfigEnvironmentTest) TestParseCmd(c *C) {
	c.Assert(parseCmd("a b"), DeepEquals, []String{"a", "b"})
	c.Assert(parseCmd("\"a b\""), DeepEquals, []String{"a b"})
	c.Assert(parseCmd("\\\"a b"), DeepEquals, []String{"\"a", "b"})
	c.Assert(parseCmd("\"\\\"a b\""), DeepEquals, []String{"\"a b"})
	c.Assert(parseCmd("\\\\"), DeepEquals, []String{"\\"})
	c.Assert(parseCmd("v1=\"a \" \"v2= b\""), DeepEquals, []String{"v1=a ", "v2= b"})
}

func Test(t *testing.T) {
	TestingT(t)
}

func init() {
	Suite(&ConfigEnvironmentTest{})
}
