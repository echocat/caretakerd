package logger

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
// The conversion pattern is "%d{YYYY-MM-DD HH:mm:ss} [%-5p]: %m%n" and the log4j environment has been set to use a PatternLayout. Then the statement will be:
// ```
// LOG debug Message 1
// LOG warn Message 2
// ```
//
// and would yield the output
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
// * ``%d[{<dateFormat>}]``: Prints out the log's creation date. Possible patterns are:
//    * Month
//       * ``M``: 1 2 ... 12
//       * ``MM``: 01 01 ... 12
//       * ``Mo``: 1st 2nd ... 12th
//       * ``MMM``: Jan Feb ... Dec
//       * ``MMMM``: January February ... December
//    * Day of Month
//       * ``D``: 1 2 ... 31
//       * ``DD``: 01 02 ... 31
//       * ``Do``: 1st 2nd ... 31st
//    * Day of Week
//       * ``ddd``: Sun Mon ... Sat
//       * ``dddd``: Sunday Monday ... Saturday
//    * Year
//       * ``YY``: 70 71 ... 12
//       * ``YYYY``: 1970 1971 ... 2012
//    * Hour
//       * ``H``: 0 1 2 ... 23
//       * ``HH``: 00 01 02 .. 23
//       * ``h``: 1 2 ... 12
//       * ``hh``: 01 02 ... 12
//    * Minute
//       * ``m``: 0 1 2 ... 59
//       * ``mm``: 00 01 02 ... 59
//    * Second
//       * ``s``: 0 1 2 ... 59
//       * ``ss``: 00 01 02 ... 59
//    * AM / PM
//       * ``A``: AM PM
//       * ``a``: am pm
//    * Timezone
//       * ``Z``: -07:00 -06:00 ... +07:00
//       * ``ZZ``: -0700 -0600 ... +0700
// * ``%m``: The log message.
// * ``%c[{<maximumNumberOfElements>}]``: Holds the logging category. Normally the instance is the name of the logger or the service. If you do not specify ``maximumNumberOfElements`` the full name is displayed. For example, if the instance is  ``%c{2}`` and the name of the category is ``a.b.c`` then the output result is ``b.c``.
// * ``%F[{<maximumNumberOfPathElements>}]``: Holds the source file that logs the instance event. If you do not specify ``maximumNumberOfPathElements`` the full file name is displayed. For example, if the instance is ``%F{2}`` and the file name is ``/a/b/c.go`` then the output result is ``b/c.go``.
// * ``%l``: Holds the source location of the log event.
// * ``%L``: Holds the line number where the log event was created.
// * ``%C[{<maximumNumberOfElements>}]``: Holds the source code package. If you do not specify ``maximumNumberOfElements`` the full name is displayed. For example, if the instance is ``%C{2}`` and the name of the package is ``a.b.c`` then the output result is ``b.c``.
// * ``%M``: Holds the method name where the log event was created.
// * ``%p``: Holds the priority or better called log level.
// * ``%P[{<subFormatPattern>}]``: Stacktrace of the location where a problem was raised that caused the instance log message.
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

// CheckedString is like String but also returns an optional error if there are any
// validation errors.
func (instance Pattern) CheckedString() (string, error) {
	return string(instance), nil
}

// Set the given string to current object from a string.
// Returns an error object if there are any problems while transforming the string.
func (instance *Pattern) Set(value string) error {
	(*instance) = Pattern(value)
	return nil
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance Pattern) MarshalYAML() (interface{}, error) {
	return string(instance), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *Pattern) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate validates action on this object and returns an error object if there are  any.
func (instance Pattern) Validate() error {
	_, err := instance.CheckedString()
	return err
}
