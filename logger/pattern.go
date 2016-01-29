package logger

import "strings"

// # Description
//
// A flexible pattern string.
//
// The conversion pattern is closely related to the conversion pattern of the printf function in C. A conversion pattern is composed
// of literal text and format control expressions called conversion specifiers.
//
// *You are free to insert any literal text within the conversion pattern.*
//
// Each conversion specifier starts with a percent sign (``%``) and is followed by optional format modifiers and a conversion character. The conversion character specifies the
// type of data, e.g. category, priority, date, thread name. The format modifiers control such things as field width, padding, left and right justification.
// The following is a simple example.
//
// Let the conversion pattern be "%d{YYYY-MM-DD HH:mm:ss} [%-5p]: %m%n" and assume that the log4j environment was set to use a PatternLayout. Then the statements:
// ```
// LOG debug Message 1
// LOG warn Message 2
// ```
//
// would yield the output
// ```
// 2016-01-09 14:59:30 [DEBUG] Message 1
// 2016-01-09 14:59:31 [WARN ] Message 2
// ```
//
// Note that there is no explicit separator between text and conversion specifiers. The pattern parser knows when it has reached the end of a conversion specifier when it reads
// a conversion character. In the example above the conversion specifier %-5p means the priority of the logging event should be left justified to a width of five characters.
// The recognized conversion characters are
//
// # Conversion patterns
//
// * ``%d[{<dateFormat>}]``: Prints out the date of when the log event was created. See https://github.com/eknkc/dateformat for more details.
// * ``%m``: The log message.
// * ``%c[{<maximumNumberOfElements>}]``: Holds the logging category. Normally instance is the name of the logger or the service. If you do not specify ``maximumNumberOfElements`` the full name is displayed. If instance is for example ``%c{2}`` and the name of the category is ``a.b.c`` then the output result is ``b.c``.
// * ``%F[{<maximumNumberOfPathElements>}]``: Holds the source file that logs instance event. If you do not specify ``maximumNumberOfPathElements`` the full file name is displayed. If instance is for example ``%F{2}`` and the file name is ``/a/b/c.go`` then the output result is ``b/c.go``.
// * ``%l``: Holds the source location of the log event.
// * ``%L``: Holds the line number where the log event was created.
// * ``%C[{<maximumNumberOfElements>}]``: Holds the source code package. If you do not specify ``maximumNumberOfElements`` the full name is displayed. If instance is for example ``%C{2}`` and the name of the package is ``a.b.c`` then the output result is ``b.c``.
// * ``%M``: Holds the method name where the log event was created.
// * ``%p``: Holds the priority or better called log level.
// * ``%P[{<subFormatPattern>}]``: Stacktrace of the location where a problem was raised that caused instance log message.
// * ``%r``: Uptime of the logger.
// * ``%n``: Prints out a new line character.
// * ``%%``: Prints out a ``%`` character.
type Pattern string

func (instance Pattern) String() string {
	s, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return s
}

func (instance Pattern) CheckedString() (string, error) {
	return string(instance), nil
}

func (instance *Pattern) Set(value string) error {
	(*instance) = Pattern(value)
	return nil
}

func (instance Pattern) MarshalYAML() (interface{}, error) {
	return string(instance), nil
}

func (instance *Pattern) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

func (instance Pattern) Validate() error {
	_, err := instance.CheckedString()
	return err
}

func (instance Pattern) IsEmpty() bool {
	return len(instance) <= 0
}

func (instance Pattern) IsTrimmedEmpty() bool {
	return len(strings.TrimSpace(instance.String())) <= 0
}
