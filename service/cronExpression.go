package service

import (
	"gopkg.in/robfig/cron.v2"
	"strings"
	"time"
)

// @serializedAs string
// # Description
//
// A cron expression represents a set of times, using 6 space-separated fields.
//
// | Field name   | Mandatory | Allowed values      | Allowed special characters     |
// | ------------ | --------- | ------------------- | ------------------------------ |
// | Seconds      | No        | ``0-59``            | ``* / , -``                    |
// | Minutes      | Yes       | ``0-59``            | ``* / , -``                    |
// | Hours        | Yes       | ``0-23``            | ``* / , -``                    |
// | Day of month | Yes       | ``1-31``            | ``* / , - ?``                  |
// | Month        | Yes       | ``1-12 or JAN-DEC`` | ``* / , -``                    |
// | Day of week  | Yes       | ``0-6 or SUN-SAT``  | ``* / , - ?``                  |
//
// > **Note:** Month and Day-of-week field values are case insensitive. ``SUN``, ``Sun``, and ``sun`` are equally accepted.
//
// # Special Characters
//
// * **Asterisk** (``*``)
// The asterisk indicates that the cron expression will match for all values of the field; e.g., using an asterisk in the
// 5th field (month) would indicate every month.
// * **Slash** (``/``)
// Slashes are used to describe increments of ranges. For example ``3-59/15`` in the 1st field (minutes) would indicate the
// 3rd minute of the hour and every 15 minutes thereafter. The form ``*\/...`` is equivalent to the form ``first-last/...``, that is, an increment
// over the largest possible range of the field. The form ``N/...`` is accepted as meaning ``N-MAX/...``, that is, starting at N, use the increment
// until the end of that specific range. It does not wrap around.
// * **Comma** (``,``)
// Commas are used to separate items of a list. For example, using ``MON,WED,FRI`` in the 5th field (day of week) would mean Mondays,
// Wednesdays and Fridays.
// * **Hyphen** (``-``)
// Hyphens are used to define ranges. For example, ``9-17`` would indicate every hour between 9am and 5pm inclusive.
// * **Question mark** (``?``)
// Question mark may be used instead of ``*`` for leaving either day-of-month or day-of-week blank.
//
// # Predefined schedules
//
// You may use one of several pre-defined schedules in place of a cron expression.
//
// | Entry                          | Description                                | Equivalent To   |
// | ------------------------------ | ------------------------------------------ | --------------- |
// | ``@yearly`` (or ``@annually)`` | Run once a year, midnight, Jan. 1st        | ``0 0 0 1 1 *`` |
// | ``@monthly``                   | Run once a month, midnight, first of month | ``0 0 0 1 * *`` |
// | ``@weekly``                    | Run once a week, midnight on Sunday        | ``0 0 0 * * 0`` |
// | ``@daily (or @midnight)``      | Run once a day, midnight                   | ``0 0 0 * * *`` |
// | ``@hourly``                    | Run once an hour, beginning of hour        | ``0 0 * * * *`` |
//
// # Intervals
//
// You may also schedule a job to execute at fixed intervals. This is supported by formatting the cron spec like this:
// ```
// @every <duration>
// ```
// where ``duration`` is a string accepted by [time.ParseDuration](http://golang.org/pkg/time/#ParseDuration).
//
// For example, ``@every 1h30m10s`` would indicate a schedule that activates every 1 hour, 30 minutes, 10 seconds.
//
// > **Hint:** The interval does not take the job runtime into account. For example, if a job takes 3 minutes to run, and it
// is scheduled to run every 5 minutes, it will have only 2 minutes of idle time between each run.
type CronExpression struct {
	spec     string
	schedule cron.Schedule
}

func NewCronExpression() CronExpression {
	return CronExpression{
		spec:     "",
		schedule: nil,
	}
}

func (instance CronExpression) String() string {
	return instance.spec
}

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
func (instance *CronExpression) Set(value string) error {
	if len(strings.TrimSpace(value)) > 0 {
		schedule, err := cron.Parse(value)
		if err != nil {
			return err
		}
		instance.schedule = schedule
		instance.spec = value
	} else {
		instance.schedule = nil
		instance.spec = ""
	}
	return nil
}

// MarshalYAML is used until yaml marshalling. Do not call directly.
func (instance CronExpression) MarshalYAML() (interface{}, error) {
	return instance.String(), nil
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call directly.
func (instance *CronExpression) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

func (instance CronExpression) IsEnabled() bool {
	return instance.schedule != nil
}

func (instance CronExpression) Next(from time.Time) *time.Time {
	if instance.IsEnabled() {
		result := instance.schedule.Next(from)
		return &result
	}
	return nil
}
