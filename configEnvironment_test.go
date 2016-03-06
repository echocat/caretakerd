package caretakerd

import (
	. "github.com/echocat/caretakerd/values"
	. "gopkg.in/check.v1"
)

type ConfigEnvironmentTest struct{}

func init() {
	Suite(&ConfigEnvironmentTest{})
}

func (s *ConfigEnvironmentTest) TestParseCmd(c *C) {
	c.Assert(parseCmd("a b"), DeepEquals, []String{"a", "b"})
	c.Assert(parseCmd("\"a b\""), DeepEquals, []String{"a b"})
	c.Assert(parseCmd("\\\"a b"), DeepEquals, []String{"\"a", "b"})
	c.Assert(parseCmd("\"\\\"a b\""), DeepEquals, []String{"\"a b"})
	c.Assert(parseCmd("\\\\"), DeepEquals, []String{"\\"})
	c.Assert(parseCmd("v1=\"a \" \"v2= b\""), DeepEquals, []String{"v1=a ", "v2= b"})
}
