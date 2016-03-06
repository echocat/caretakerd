package access

import (
	"github.com/echocat/caretakerd/errors"
	. "gopkg.in/check.v1"
	"os"
	"reflect"
)

type FilePermissionTest struct{}

func init() {
	Suite(&FilePermissionTest{})
}

func (s *FilePermissionTest) TestString(c *C) {
	c.Assert(FilePermission(0600).String(), Equals, "0600")
	c.Assert(FilePermission(0).String(), Equals, "0000")
}

func (s *FilePermissionTest) TestSet(c *C) {
	actual := FilePermission(0)
	c.Assert(actual, Equals, FilePermission(0))

	c.Assert(actual.Set("0600"), IsNil)
	c.Assert(actual, Equals, FilePermission(0600))

	c.Assert(actual.Set("666600"), ErrorMatches, "Illegal file permission format.*")
	c.Assert(actual.Set("xx"), ErrorMatches, "Illegal file permission format.*")
}

func (s *FilePermissionTest) TestMarshalYAML(c *C) {
	actual := FilePermission(0600)
	pb, err := actual.MarshalYAML()
	c.Assert(err, IsNil)
	c.Assert(pb.(string), Equals, "0600")
}

func (s *FilePermissionTest) TestUnmarshalYAML(c *C) {
	actual := FilePermission(0)
	c.Assert(actual.UnmarshalYAML(func(v interface{}) error {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf("0600"))
		return nil
	}), IsNil)
	c.Assert(actual, Equals, FilePermission(0600))
}

func (s *FilePermissionTest) TestUnmarshalYAMLWithProblems(c *C) {
	actual := FilePermission(0)
	c.Assert(actual.UnmarshalYAML(func(v interface{}) error {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf("0000600"))
		return nil
	}), ErrorMatches, "Illegal file permission format: 0000600")
	c.Assert(actual, Equals, FilePermission(0))

	c.Assert(actual.UnmarshalYAML(func(v interface{}) error {
		return errors.New("Foo")
	}), ErrorMatches, "Foo")
	c.Assert(actual, Equals, FilePermission(0))
}

func (s *FilePermissionTest) TestMarshalJSON(c *C) {
	actual := FilePermission(0600)
	pb, err := actual.MarshalJSON()
	c.Assert(err, IsNil)
	c.Assert(string(pb), Equals, "\"0600\"")
}

func (s *FilePermissionTest) TestUnmarshalJSON(c *C) {
	actual := FilePermission(0)
	c.Assert(actual.UnmarshalJSON([]byte("\"0600\"")), IsNil)
	c.Assert(actual, Equals, FilePermission(0600))
}

func (s *FilePermissionTest) TestUnmarshalJSONWithProblems(c *C) {
	actual := FilePermission(0)
	c.Assert(actual.UnmarshalJSON([]byte("\"0000600\"")), ErrorMatches, "Illegal file permission format: 0000600")
	c.Assert(actual, Equals, FilePermission(0))

	c.Assert(actual.UnmarshalJSON([]byte("0000600")), ErrorMatches, "invalid character '0' after top-level value")
	c.Assert(actual, Equals, FilePermission(0))
}

func (s *FilePermissionTest) TestValidate(c *C) {
	c.Assert(FilePermission(0).Validate(), IsNil)
}

func (s *FilePermissionTest) TestThisOrDefault(c *C) {
	c.Assert(FilePermission(0).ThisOrDefault(), Equals, DefaultFilePermission())
	c.Assert(FilePermission(0644).ThisOrDefault(), Equals, FilePermission(0644))
}

func (s *FilePermissionTest) TestAsFileMode(c *C) {
	c.Assert(FilePermission(0).AsFileMode(), Equals, os.FileMode(0))
	c.Assert(FilePermission(0600).AsFileMode(), Equals, os.FileMode(0600))
}
