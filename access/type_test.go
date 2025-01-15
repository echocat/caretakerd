package access

import (
	"reflect"

	. "gopkg.in/check.v1"

	"github.com/echocat/caretakerd/errors"
	. "github.com/echocat/caretakerd/testing"
)

type TypeTest struct{}

func init() {
	Suite(&TypeTest{})
}

func (s *TypeTest) TestString(c *C) {
	c.Assert(None.String(), Equals, "none")
	c.Assert(Trusted.String(), Equals, "trusted")
	c.Assert(GenerateToEnvironment.String(), Equals, "generateToEnvironment")
	c.Assert(GenerateToFile.String(), Equals, "generateToFile")
}

func (s *TypeTest) TestStringPanic(c *C) {
	c.Assert(func() {
		_ = Type(-1).String()
	}, ThrowsPanicThatMatches, "illegal access type: -1")
}

func (s *TypeTest) TestCheckedString(c *C) {
	for _, t := range AllTypes {
		r, err := t.CheckedString()
		c.Assert(err, IsNil)
		c.Assert(r, Not(IsEmpty))
	}
}

func (s *TypeTest) TestCheckedStringErrors(c *C) {
	r, err := Type(-1).CheckedString()
	c.Assert(err, ErrorMatches, "illegal access type: -1")
	c.Assert(r, Equals, "")
}

func (s *TypeTest) TestSet(c *C) {
	actual := Type(-1)
	c.Assert(actual.Set("1"), IsNil)
	c.Assert(actual, Equals, Trusted)

	c.Assert(actual.Set("generateToEnvironment"), IsNil)
	c.Assert(actual, Equals, GenerateToEnvironment)

	c.Assert(actual.Set("xxx"), ErrorMatches, "illegal access type: xxx")
	c.Assert(actual, Equals, GenerateToEnvironment)

	c.Assert(actual.Set("666"), ErrorMatches, "illegal access type: 666")
	c.Assert(actual, Equals, GenerateToEnvironment)
}

func (s *TypeTest) TestMarshalYAML(c *C) {
	actual := Trusted
	pb, err := actual.MarshalYAML()
	c.Assert(err, IsNil)
	c.Assert(pb.(string), Equals, "trusted")
}

func (s *TypeTest) TestUnmarshalYAML(c *C) {
	actual := Type(-1)
	c.Assert(actual.UnmarshalYAML(func(v interface{}) error {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf("trusted"))
		return nil
	}), IsNil)
	c.Assert(actual, Equals, Trusted)
}

func (s *TypeTest) TestUnmarshalYAMLWithProblems(c *C) {
	actual := Type(-1)
	c.Assert(actual.UnmarshalYAML(func(v interface{}) error {
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf("foobar"))
		return nil
	}), ErrorMatches, "illegal access type: foobar")
	c.Assert(actual, Equals, Type(-1))

	c.Assert(actual.UnmarshalYAML(func(v interface{}) error {
		return errors.New("Foo")
	}), ErrorMatches, "Foo")
	c.Assert(actual, Equals, Type(-1))
}

func (s *TypeTest) TestMarshalJSON(c *C) {
	actual := Trusted
	pb, err := actual.MarshalJSON()
	c.Assert(err, IsNil)
	c.Assert(string(pb), Equals, "\"trusted\"")
}

func (s *TypeTest) TestMarshalJSONWithProblems(c *C) {
	actual := Type(-1)
	pb, err := actual.MarshalJSON()
	c.Assert(err, ErrorMatches, "illegal access type: -1")
	c.Assert(string(pb), Equals, "")
}

func (s *TypeTest) TestUnmarshalJSON(c *C) {
	actual := Type(-1)
	c.Assert(actual.UnmarshalJSON([]byte("\"trusted\"")), IsNil)
	c.Assert(actual, Equals, Trusted)
}

func (s *TypeTest) TestUnmarshalJSONWithProblems(c *C) {
	actual := Trusted
	c.Assert(actual.UnmarshalJSON([]byte("\"foobar\"")), ErrorMatches, "illegal access type: foobar")
	c.Assert(actual, Equals, Trusted)

	c.Assert(actual.UnmarshalJSON([]byte("0000")), ErrorMatches, "invalid character '0' after top-level value")
	c.Assert(actual, Equals, Trusted)
}

func (s *TypeTest) TestValidate(c *C) {
	c.Assert(Trusted.Validate(), IsNil)
	c.Assert(Type(-1).Validate(), ErrorMatches, "illegal access type: -1")
}

func (s *TypeTest) TestIsTakingFilename(c *C) {
	for _, t := range AllTypes {
		c.Assert(t.IsTakingFilename(), Equals, t == GenerateToFile)
	}
}

func (s *TypeTest) TestIsTakingFilePermission(c *C) {
	for _, t := range AllTypes {
		c.Assert(t.IsTakingFilePermission(), Equals, t == GenerateToFile)
	}
}

func (s *TypeTest) TestIsTakingFileUser(c *C) {
	for _, t := range AllTypes {
		c.Assert(t.IsTakingFileUser(), Equals, t == GenerateToFile)
	}
}

func (s *TypeTest) TestIsTakingFileGroup(c *C) {
	for _, t := range AllTypes {
		c.Assert(t.IsTakingFileGroup(), Equals, t == GenerateToFile)
	}
}

func (s *TypeTest) TestIsGenerating(c *C) {
	for _, t := range AllTypes {
		c.Assert(t.IsGenerating(), Equals, t == GenerateToFile || t == GenerateToEnvironment)
	}
}
