package service

import (
    "gopkg.in/robfig/cron.v2"
    "strings"
    "time"
)

type CronExpression struct {
    spec     string
    schedule cron.Schedule
}

func NewCronExpression() CronExpression {
    return CronExpression{
        spec: "",
        schedule: nil,
    }
}

func (this CronExpression) String() string {
    return this.spec
}

func (this *CronExpression) Set(value string) error {
    if len(strings.TrimSpace(value)) > 0 {
        schedule, err := cron.Parse(value)
        if err != nil {
            return err
        }
        this.schedule = schedule
        this.spec = value
    } else {
        this.schedule = nil
        this.spec = ""
    }
    return nil
}

func (this CronExpression) MarshalYAML() (interface{}, error) {
    return this.String(), nil
}

func (this *CronExpression) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value string
    if err := unmarshal(&value); err != nil {
        return err
    }
    return this.Set(value)
}

func (this CronExpression) IsEnabled() bool {
    return this.schedule != nil
}

func (this CronExpression) Next(from time.Time) *time.Time {
    if this.IsEnabled() {
        result := this.schedule.Next(from)
        return &result
    }
    return nil
}
