package service

import (
	"gopkg.in/robfig/cron.v2"
	"strings"
	"time"
)

// @serializedAs string
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

func (instance CronExpression) MarshalYAML() (interface{}, error) {
	return instance.String(), nil
}

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
