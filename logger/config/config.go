package config

import (
    . "github.com/echocat/caretakerd/logger/level"
    . "github.com/echocat/caretakerd/values"
)

var defaults = map[string]interface{} {
    "Level": Info,
    "StdoutLevel": Info,
    "StderrLevel": Error,
    "Filename": String("console"),
    "MaxSizeInMb": NonNegativeInteger(500),
    "MaxBackups": NonNegativeInteger(5),
    "MaxAgeInDays": NonNegativeInteger(1),
    "Pattern": String("%d{YYYY-MM-DD HH:mm:ss} [%-5.5p] [%c] %m%n%P{%m}"),
}

type Config struct {
    Level        Level `json:"level" yaml:"level"`
    StdoutLevel  Level `json:"stdoutLevel" yaml:"stdoutLevel"`
    StderrLevel  Level `json:"stderrLevel" yaml:"stderrLevel"`
    Filename     String `json:"filename" yaml:"filename"`
    MaxSizeInMb  NonNegativeInteger `json:"maxSizeInMb" yaml:"maxSizeInMb"`
    MaxBackups   NonNegativeInteger `json:"maxBackups" yaml:"maxBackups"`
    MaxAgeInDays NonNegativeInteger `json:"maxAgeInDays" yaml:"maxAgeInDays"`
    Pattern      String  `json:"pattern" yaml:"pattern"`
}

func NewLoggerConfig() Config {
    result := Config{}
    result.init()
    return result
}

func (this Config) Validate() error {
    err := this.StdoutLevel.Validate()
    if err == nil {
        err = this.StderrLevel.Validate()
    }
    return err
}

func (this *Config) init() {
    SetDefaultsTo(defaults, this)
}

func (this *Config) BeforeUnmarshalYAML() error {
    this.init()
    return nil
}
