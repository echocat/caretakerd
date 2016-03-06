package access

import (
	"github.com/echocat/caretakerd/errors"
	. "github.com/echocat/caretakerd/testing"
	. "gopkg.in/check.v1"
	"reflect"
)

type PermissionTest struct{}

func init() {
	Suite(&PermissionTest{})
}

func (s *PermissionTest) TestString(c *C) {
	c.Assert(Forbidden.String(), Equals, "forbidden")
	c.Assert(ReadOnly.String(), Equals, "readOnly")
	c.Assert(ReadWrite.String(), Equals, "readWrite")
}

func (s *PermissionTest) TestStringPanic(c *C) {
	c.Assert(func() {
		Permission(-1).String()
	}, ThrowsPanicThatMatches, "Illegal permission: -1")
}

func (s *PermissionTest) TestCheckedString(c *C) {
	for _, permission := range AllPermissions {
		r, err := permission.CheckedString()
		c.Assert(err, IsNil)
		c.Assert(r, Not(IsEmpty))
	}
}

func (s *PermissionTest) TestCheckedStringErrors(c *C) {
	r, err := Permission(-1).CheckedString()
	c.Assert(err, ErrorMatches, "Illegal permission: -1")
	c.Assert(r, Equals, "")
}

func (s *PermissionTest) TestSet(c *C) {
	actual := Permission(-1)
	c.Assert(actual.Set("1"), IsNil)
	c.Assert(actual, Equals, ReadOnly)

	c.Assert(actual.Set("readWrite"), IsNil)
	c.Assert(actual, Equals, ReadWrite)

	c.Assert(actual.Set("xxx"), ErrorMatches, "Illegal permission: xxx")
	c.Assert(actual, Equals, ReadWrite)

	c.Assert(actual.Set("666"), ErrorMatches, "Illegal permission: 666")
	c.Assert(actual, Equals, ReadWrite)
}

func (s *PermissionTest) TestMarshalYAML(c *C) {
	actual := ReadOnly
	pb, err := actual.MarshalYAML()
	c.Assert(err, IsNil)
	c.Assert(pb.(string), Equals, "readOnly")
}

func (s *PermissionTest) TestUnmarshalYAML(c *C) {
	actual := Permission(-1)
	c.Assert(actual.UnmarshalYAML(func(v interface{}) error {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf("readOnly"))
		return nil
	}), IsNil)
	c.Assert(actual, Equals, ReadOnly)
}

func (s *PermissionTest) TestUnmarshalYAMLWithProblems(c *C) {
	actual := Permission(-1)
	c.Assert(actual.UnmarshalYAML(func(v interface{}) error {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf("foobar"))
		return nil
	}), ErrorMatches, "Illegal permission: foobar")
	c.Assert(actual, Equals, Permission(-1))

	c.Assert(actual.UnmarshalYAML(func(v interface{}) error {
		return errors.New("Foo")
	}), ErrorMatches, "Foo")
	c.Assert(actual, Equals, Permission(-1))
}

func (s *PermissionTest) TestMarshalJSON(c *C) {
	actual := ReadOnly
	pb, err := actual.MarshalJSON()
	c.Assert(err, IsNil)
	c.Assert(string(pb), Equals, "\"readOnly\"")
}

func (s *PermissionTest) TestMarshalJSONWithProblems(c *C) {
	actual := Permission(-1)
	pb, err := actual.MarshalJSON()
	c.Assert(err, ErrorMatches, "Illegal permission: -1")
	c.Assert(string(pb), Equals, "")
}

func (s *PermissionTest) TestUnmarshalJSON(c *C) {
	actual := Permission(-1)
	c.Assert(actual.UnmarshalJSON([]byte("\"readOnly\"")), IsNil)
	c.Assert(actual, Equals, ReadOnly)
}

func (s *PermissionTest) TestUnmarshalJSONWithProblems(c *C) {
	actual := ReadOnly
	c.Assert(actual.UnmarshalJSON([]byte("\"foobar\"")), ErrorMatches, "Illegal permission: foobar")
	c.Assert(actual, Equals, ReadOnly)

	c.Assert(actual.UnmarshalJSON([]byte("0000")), ErrorMatches, "invalid character '0' after top-level value")
	c.Assert(actual, Equals, ReadOnly)
}

func (s *PermissionTest) TestValidate(c *C) {
	c.Assert(ReadOnly.Validate(), IsNil)
	c.Assert(Permission(-1).Validate(), ErrorMatches, "Illegal permission: -1")
}
