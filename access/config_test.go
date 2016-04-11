package access

import (
	. "github.com/echocat/caretakerd/testing"
	"github.com/echocat/caretakerd/values"
	. "gopkg.in/check.v1"
)

type ConfigTest struct{}

func init() {
	Suite(&ConfigTest{})
}

func (s *ConfigTest) TestNewNoneConfig(c *C) {
	actual := NewNoneConfig()
	c.Assert(actual.Type, Equals, None)
	c.Assert(actual.Permission, Equals, Forbidden)
	c.Assert(actual.PemFile, IsEmpty)
	c.Assert(actual.PemFilePermission, Equals, FilePermission(0))
	c.Assert(actual.PemFileUser, IsEmpty)
}

func (s *ConfigTest) TestNewTrustedConfig(c *C) {
	actual := NewTrustedConfig(ReadOnly)
	c.Assert(actual.Type, Equals, Trusted)
	c.Assert(actual.Permission, Equals, ReadOnly)
	c.Assert(actual.PemFile, IsEmpty)
	c.Assert(actual.PemFilePermission, Equals, FilePermission(0))
	c.Assert(actual.PemFileUser, IsEmpty)
}

func (s *ConfigTest) TestNewGenerateToEnvironmentConfig(c *C) {
	actual := NewGenerateToEnvironmentConfig(ReadOnly)
	c.Assert(actual.Type, Equals, GenerateToEnvironment)
	c.Assert(actual.Permission, Equals, ReadOnly)
	c.Assert(actual.PemFile, IsEmpty)
	c.Assert(actual.PemFilePermission, Equals, FilePermission(0))
	c.Assert(actual.PemFileUser, IsEmpty)
}

func (s *ConfigTest) TestNewGenerateToFileConfig(c *C) {
	actual := NewGenerateToFileConfig(ReadOnly, values.String("foo/bar.pem"))
	c.Assert(actual.Type, Equals, GenerateToFile)
	c.Assert(actual.Permission, Equals, ReadOnly)
	c.Assert(actual.PemFile, Equals, values.String("foo/bar.pem"))
	c.Assert(actual.PemFilePermission, Equals, DefaultFilePermission())
	c.Assert(actual.PemFileUser, IsEmpty)
}

func (s *ConfigTest) TestValidateWrongType(c *C) {
	actual := Config{
		Type: Type(66),
	}
	c.Assert(actual.Validate(), ErrorMatches, "Illegal access type: 66")
	actual.Type = Trusted
	c.Assert(actual.Validate(), IsNil)
}

func (s *ConfigTest) TestValidateWrongPermission(c *C) {
	actual := Config{
		Type:       Trusted,
		Permission: Permission(66),
	}
	c.Assert(actual.Validate(), ErrorMatches, "Illegal permission: 66")
	actual.Permission = ReadOnly
	c.Assert(actual.Validate(), IsNil)
}

func (s *ConfigTest) TestValidateRequiredPemFile(c *C) {
	for _, t := range AllTypes {
		if t.IsTakingFilename() {
			actual := Config{
				Type:       t,
				Permission: ReadOnly,
				PemFile:    values.String(""),
			}
			c.Assert(actual.Validate(), ErrorMatches, "There is no pemFile set for type "+t.String()+".")
			actual.PemFile = values.String("foo/bar.pem")
			c.Assert(actual.Validate(), IsNil)
		}
	}
}

func (s *ConfigTest) TestValidateNotRequiredPemFile(c *C) {
	for _, t := range AllTypes {
		if !t.IsTakingFilename() {
			actual := Config{
				Type:       t,
				Permission: ReadOnly,
				PemFile:    values.String(""),
			}
			c.Assert(actual.Validate(), IsNil)
			actual.PemFile = values.String("foo/bar.pem")
			c.Assert(actual.Validate(), IsNil)
		}
	}
}

func (s *ConfigTest) TestValidateAllowedPemFileUser(c *C) {
	for _, t := range AllTypes {
		if t.IsTakingFileUser() {
			actual := Config{
				Type:       t,
				Permission: ReadOnly,
				PemFile:    values.String("foo/bar.pem"),
			}
			actual.PemFileUser = values.String("")
			c.Assert(actual.Validate(), IsNil)
			actual.PemFileUser = values.String("foo")
			c.Assert(actual.Validate(), IsNil)
		}
	}
}

func (s *ConfigTest) TestValidateNotAllowedPemFileUser(c *C) {
	for _, t := range AllTypes {
		if !t.IsTakingFileUser() {
			actual := Config{
				Type:       t,
				Permission: ReadOnly,
				PemFile:    values.String(""),
			}
			actual.PemFileUser = values.String("")
			c.Assert(actual.Validate(), IsNil)
			actual.PemFileUser = values.String("foo")
			c.Assert(actual.Validate(), ErrorMatches, "There is no pemFileUser allowed for type "+t.String()+".")
		}
	}
}

func (s *ConfigTest) TestValidateAllowedPemFilePermission(c *C) {
	for _, t := range AllTypes {
		if t.IsTakingFilePermission() {
			actual := Config{
				Type:       t,
				Permission: ReadOnly,
				PemFile:    values.String("foo/bar.pem"),
			}
			actual.PemFilePermission = FilePermission(0)
			c.Assert(actual.Validate(), IsNil)
			actual.PemFilePermission = FilePermission(0600)
			c.Assert(actual.Validate(), IsNil)
		}
	}
}

func (s *ConfigTest) TestValidateNotAllowedFilePermission(c *C) {
	for _, t := range AllTypes {
		if !t.IsTakingFilePermission() {
			actual := Config{
				Type:       t,
				Permission: ReadOnly,
				PemFile:    values.String(""),
			}
			actual.PemFilePermission = FilePermission(0)
			c.Assert(actual.Validate(), IsNil)
			actual.PemFilePermission = DefaultFilePermission()
			c.Assert(actual.Validate(), IsNil)
			actual.PemFilePermission = FilePermission(0611)
			c.Assert(actual.Validate(), ErrorMatches, "There is no pemFilePermission allowed for type "+t.String()+".")
		}
	}
}
