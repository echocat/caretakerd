package caretakerd

import (
	"os"
	"strings"
	"github.com/echocat/caretakerd/panics"
	"github.com/echocat/caretakerd/service"
	"github.com/echocat/caretakerd/errors"
	. "github.com/echocat/caretakerd/values"
)

const (
	envPrefix = "CTD."
)

var globalEnvKeyToHandler map[string]func(conf *Config, value string) error = map[string]func(*Config, string) error{
	// logger.config
	"LOG_LEVEL": handleGlobalLogLevelEnv,
	"LOG_STDOUT_LEVEL": handleGlobalLogStdoutLevelEnv,
	"LOG_STDERR_LEVEL": handleGlobalLogStderrLevelEnv,
	"LOG_FILE": handleGlobalLogFilenameEnv,
	"LOG_FILE_NAME": handleGlobalLogFilenameEnv,
	"LOG_MAX_SIZE": handleGlobalLogMaxSizeInMbEnv,
	"LOG_MAX_SIZE_IN_MB": handleGlobalLogMaxSizeInMbEnv,
	"LOG_MAX_BACKUPS": handleGlobalLogMaxBackupsEnv,
	"LOG_MAX_AGE": handleGlobalLogMaxAgeInDaysEnv,
	"LOG_MAX_AGE_IN_DAYS": handleGlobalLogMaxAgeInDaysEnv,
}

var serviceEnvKeyToFunction map[string]func(conf *service.Config, value string) error = map[string]func(*service.Config, string) error{
	// service.config
	"CMD": handleServiceCommandEnv,
	"COMMAND": handleServiceCommandEnv,
	"TYPE": handleServiceTypeEnv,
	"START_DELAY": handleServiceStartDelayInSecondsEnv,
	"START_DELAY_IN_SECONDS": handleServiceStartDelayInSecondsEnv,
	"RESTART_DELAY": handleServiceRestartDelayInSecondsEnv,
	"RESTART_DELAY_IN_SECONDS": handleServiceRestartDelayInSecondsEnv,
	"SUCCESS": handleServiceSuccessExitCodesEnv,
	"EXIT_CODE": handleServiceSuccessExitCodesEnv,
	"EXIT_CODES": handleServiceSuccessExitCodesEnv,
	"SUCCESS_EXIT_CODE": handleServiceSuccessExitCodesEnv,
	"SUCCESS_EXIT_CODES": handleServiceSuccessExitCodesEnv,
	"STOP_WAIT": handleServiceStopWaitInSecondsEnv,
	"STOP_WAIT_IN_SECONDS": handleServiceStopWaitInSecondsEnv,
	"USER": handleServiceUserEnv,
	"DIR": handleServiceDirectoryEnv,
	"DIRECORY": handleServiceDirectoryEnv,
	"RESTART": handleServiceAutoRestartEnv,
	"AUTO_RESTART": handleServiceAutoRestartEnv,
	"INHERIT_ENV": handleServiceInheritEnvironmentEnv,
	"INHERIT_ENVIRONMENT": handleServiceInheritEnvironmentEnv,
	// logger.config
	"LOG_LEVEL": handleServiceLogLevelEnv,
	"LOG_STDOUT_LEVEL": handleServiceLogStdoutLevelEnv,
	"LOG_STDERR_LEVEL": handleServiceLogStderrLevelEnv,
	"LOG_FILE": handleServiceLogFilenameEnv,
	"LOG_FILE_NAME": handleServiceLogFilenameEnv,
	"LOG_MAX_SIZE": handleServiceLogMaxSizeInMbEnv,
	"LOG_MAX_SIZE_IN_MB": handleServiceLogMaxSizeInMbEnv,
	"LOG_MAX_BACKUPS": handleServiceLogMaxBackupsEnv,
	"LOG_MAX_AGE": handleServiceLogMaxAgeInDaysEnv,
	"LOG_MAX_AGE_IN_DAYS": handleServiceLogMaxAgeInDaysEnv,
}

var serviceSubEnvKeyToFunction map[string]func(conf *service.Config, key string, value string) error = map[string]func(*service.Config, string, string) error{
	// service.config
	"ENV": handleEnvironmentEnv,
	"ENVIRONMENT": handleEnvironmentEnv,
}

type Appendable interface {
	Append(value string) error
}

type Putable interface {
	Put(key string, value string) error
}

func (i Config) EnrichFromEnvironment() Config {
	result := &i
	for _, plainEnviron := range os.Environ() {
		environ := strings.SplitN(plainEnviron, "=", 2)
		var err error
		if len(environ) > 1 {
			err = i.handleUnknownEnv(plainEnviron, environ[0], environ[1])
		} else {
			err = i.handleUnknownEnv(plainEnviron, environ[0], "")
		}
		if err != nil {
			panics.New("Could not handle environment variable '%s'. Got: %s", plainEnviron, err.Error()).CausedBy(err).Throw()
		}
	}
	return *result
}

func (i *Config) handleUnknownEnv(full string, key string, value string) error {
	prefixLength := len(envPrefix)
	var err error
	if strings.HasPrefix(strings.ToUpper(key), envPrefix) && len(key) > prefixLength {
		err = i.handleEnv(full, key[prefixLength:], value)
	} else {
		err = nil
	}
	return err
}

func (i *Config) handleEnv(full string, key string, value string) error {
	if handler, ok := globalEnvKeyToHandler[key]; ok {
		return handler(i, value)
	} else {
		parts := strings.SplitN(key, ".", 3)
		var err error
		if len(parts) == 2 {
			err = i.handleServiceEnv(full, parts[0], parts[1], value)
		} else if len(parts) == 3 {
			err = i.handleServiceMapEnv(full, parts[0], parts[1], parts[2], value)
		} else {
			err = errors.New("Illegal environment variable found: '%s'. Unexpected number of '.' chracters.", full)
		}
		return err
	}
}

func handleGlobalLogLevelEnv(conf *Config, value string) error {
	return conf.Logger.Level.Set(value)
}

func handleGlobalLogStdoutLevelEnv(conf *Config, value string) error {
	return conf.Logger.StdoutLevel.Set(value)
}

func handleGlobalLogStderrLevelEnv(conf *Config, value string) error {
	return conf.Logger.StderrLevel.Set(value)
}

func handleGlobalLogFilenameEnv(conf *Config, value string) error {
	return conf.Logger.Filename.Set(value)
}

func handleGlobalLogMaxSizeInMbEnv(conf *Config, value string) error {
	return conf.Logger.MaxSizeInMb.Set(value)
}

func handleGlobalLogMaxBackupsEnv(conf *Config, value string) error {
	return conf.Logger.MaxBackups.Set(value)
}

func handleGlobalLogMaxAgeInDaysEnv(conf *Config, value string) error {
	return conf.Logger.MaxAgeInDays.Set(value)
}

func handleServiceCommandEnv(conf *service.Config, value string) error {
	conf.Command = append(conf.Command, parseCmd(value)...)
	return nil
}

func handleServiceTypeEnv(conf *service.Config, value string) error {
	return conf.Type.Set(value)
}

func handleServiceStartDelayInSecondsEnv(conf *service.Config, value string) error {
	return conf.StartDelayInSeconds.Set(value)
}

func handleServiceRestartDelayInSecondsEnv(conf *service.Config, value string) error {
	return conf.RestartDelayInSeconds.Set(value)
}

func handleServiceSuccessExitCodesEnv(conf *service.Config, value string) error {
	return conf.SuccessExitCodes.Set(value)
}

func handleServiceStopWaitInSecondsEnv(conf *service.Config, value string) error {
	return conf.StopWaitInSeconds.Set(value)
}

func handleServiceUserEnv(conf *service.Config, value string) error {
	return conf.User.Set(value)
}

func handleServiceDirectoryEnv(conf *service.Config, value string) error {
	return conf.Directory.Set(value)
}

func handleServiceAutoRestartEnv(conf *service.Config, value string) error {
	return conf.AutoRestart.Set(value)
}

func handleServiceLogLevelEnv(conf *service.Config, value string) error {
	return conf.Logger.Level.Set(value)
}

func handleServiceLogStdoutLevelEnv(conf *service.Config, value string) error {
	return conf.Logger.StdoutLevel.Set(value)
}

func handleServiceLogStderrLevelEnv(conf *service.Config, value string) error {
	return conf.Logger.StderrLevel.Set(value)
}

func handleServiceLogFilenameEnv(conf *service.Config, value string) error {
	return conf.Logger.Filename.Set(value)
}

func handleServiceLogMaxSizeInMbEnv(conf *service.Config, value string) error {
	return conf.Logger.MaxSizeInMb.Set(value)
}

func handleServiceLogMaxBackupsEnv(conf *service.Config, value string) error {
	return conf.Logger.MaxBackups.Set(value)
}

func handleServiceLogMaxAgeInDaysEnv(conf *service.Config, value string) error {
	return conf.Logger.MaxAgeInDays.Set(value)
}

func handleServiceInheritEnvironmentEnv(conf *service.Config, value string) error {
	return conf.InheritEnvironment.Set(value)
}

func (i *Config) handleServiceEnv(full string, serviceName string, key string, value string) error {
	targetKey := strings.ToUpper(key)
	if handler, ok := serviceEnvKeyToFunction[targetKey]; ok {
		return i.Services.Configure(serviceName, value, handler)
	} else {
		return errors.New("Unknown configuration type '%s' for service '%s'.", key, serviceName)
	}
}

func handleEnvironmentEnv(conf *service.Config, key string, value string) error {
	return conf.Environment.Put(key, value)
}

func (i *Config) handleServiceMapEnv(full string, serviceName string, key string, subKey string, value string) error {
	targetKey := strings.ToUpper(key)
	if handler, ok := serviceSubEnvKeyToFunction[targetKey]; ok {
		return i.Services.ConfigureSub(serviceName, subKey, value, handler)
	} else {
		return errors.New("Unknown configuration type '%s' for service '%s'.", key, serviceName)
	}
}

func parseCmd(in string) []String {
	result := []String{}
	inEscape := false
	inQuotes := false
	buf := ""
	for i := 0; i < len(in); i++ {
		c := in[i]
		if inEscape {
			buf += string(c)
			inEscape = false
		} else if c == '\\' {
			inEscape = true
		} else if inQuotes {
			if c == '"' {
				inQuotes = false
			} else {
				buf += string(c)
			}
		} else if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			result = append(result, String(buf))
			buf = ""
		} else if c == '"' {
			inQuotes = true
		} else {
			buf += string(c)
		}
	}
	if len(buf) > 0 {
		result = append(result, String(buf))
	}
	return result
}
