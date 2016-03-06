package testing

import (
	"errors"
	"github.com/echocat/caretakerd/values"
	. "gopkg.in/check.v1"
	"time"
)

type CheckersTest struct{}

func init() {
	Suite(&CheckersTest{})
}

type enclosesString struct {
	string
}

func (instance enclosesString) String() string {
	return instance.string
}

func (s *CheckersTest) TestThrowsPanicThatMatches(c *C) {
	c.Assert(func() {
		panic(errors.New("foo123"))
	}, ThrowsPanicThatMatches, "foo1.3")
	c.Assert(func() {
		panic(enclosesString{string: "foo123"})
	}, ThrowsPanicThatMatches, "foo1.3")
	c.Assert(func() {
		panic("foo123")
	}, ThrowsPanicThatMatches, "foo1.3")
}

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

func (s *CheckersTest) TestIsLessThan(c *C) {
	c.Assert(22, IsLessThan, 66)
	c.Assert(22, Not(IsLessThan), 11)
	c.Assert(22, Not(IsLessThan), 22)
	c.Assert(22.5, IsLessThan, 66.5)
	c.Assert(22.5, Not(IsLessThan), 11.5)
	c.Assert(22.5, Not(IsLessThan), 22.5)
	c.Assert(time.Duration(22), IsLessThan, time.Duration(66))

	result, message := IsLessThan.Check([]interface{}{22.5, 11}, []string{"obtained", "compareTo"})
	c.Assert(result, Equals, false)
	c.Assert(message, Equals, "'obtained' type not equal to the type to 'compareTo' type.")
}

func (s *CheckersTest) TestIsLessThanOrEqual(c *C) {
	c.Assert(22, IsLessThanOrEqualTo, 66)
	c.Assert(22, IsLessThanOrEqualTo, 22)
	c.Assert(22, Not(IsLessThanOrEqualTo), 11)
	c.Assert(22.5, IsLessThanOrEqualTo, 66.5)
	c.Assert(22.5, IsLessThanOrEqualTo, 22.5)
	c.Assert(22.5, Not(IsLessThanOrEqualTo), 11.5)
	c.Assert(time.Duration(22), IsLessThanOrEqualTo, time.Duration(66))
	c.Assert(time.Duration(22), IsLessThanOrEqualTo, time.Duration(22))

	result, message := IsLessThanOrEqualTo.Check([]interface{}{22.5, 11}, []string{"obtained", "compareTo"})
	c.Assert(result, Equals, false)
	c.Assert(message, Equals, "'obtained' type not equal to the type to 'compareTo' type.")
}

func (s *CheckersTest) TestIsLargerThan(c *C) {
	c.Assert(22, IsLargerThan, 11)
	c.Assert(22, Not(IsLargerThan), 66)
	c.Assert(22, Not(IsLargerThan), 22)
	c.Assert(66.5, IsLargerThan, 22.5)
	c.Assert(11.5, Not(IsLargerThan), 22.5)
	c.Assert(22.5, Not(IsLargerThan), 22.5)
	c.Assert(time.Duration(22), IsLargerThan, time.Duration(11))

	result, message := IsLargerThan.Check([]interface{}{22.5, 11}, []string{"obtained", "compareTo"})
	c.Assert(result, Equals, false)
	c.Assert(message, Equals, "'obtained' type not equal to the type to 'compareTo' type.")
}

func (s *CheckersTest) TestIsLargerThanOrEqualTo(c *C) {
	c.Assert(22, IsLargerThanOrEqualTo, 11)
	c.Assert(22, Not(IsLargerThanOrEqualTo), 66)
	c.Assert(22, IsLargerThanOrEqualTo, 22)
	c.Assert(66.5, IsLargerThanOrEqualTo, 22.5)
	c.Assert(11.5, Not(IsLargerThanOrEqualTo), 22.5)
	c.Assert(22.5, IsLargerThanOrEqualTo, 22.5)
	c.Assert(time.Duration(22), IsLargerThanOrEqualTo, time.Duration(11))
	c.Assert(time.Duration(22), IsLargerThanOrEqualTo, time.Duration(22))

	result, message := IsLargerThanOrEqualTo.Check([]interface{}{22.5, 11}, []string{"obtained", "compareTo"})
	c.Assert(result, Equals, false)
	c.Assert(message, Equals, "'obtained' type not equal to the type to 'compareTo' type.")
}
