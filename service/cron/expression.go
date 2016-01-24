package cron

import (
    "gopkg.in/robfig/cron.v2"
    "strings"
    "time"
)

type Expression struct {
    spec     string
    schedule cron.Schedule
}

func NewCronExpression() Expression {
    return Expression{
        spec: "",
        schedule: nil,
    }
}

func (this Expression) String() string {
    return this.spec
}

func (this *Expression) Set(value string) error {
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

func (this Expression) MarshalYAML() (interface{}, error) {
    return this.String(), nil
}

func (this *Expression) UnmarshalYAML(unmarshal func(interface{}) error) error {
    var value string
    if err := unmarshal(&value); err != nil {
        return err
    }
    return this.Set(value)
}

func (this Expression) IsEnabled() bool {
    return this.schedule != nil
}

func (this Expression) Next(from time.Time) *time.Time {
    if this.IsEnabled() {
        result := this.schedule.Next(from)
        return &result
    }
    return nil
}
