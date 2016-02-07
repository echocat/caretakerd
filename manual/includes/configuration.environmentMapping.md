# Environment mapping {#environmentMapping}

There are several environment variables mapped to configuration parameters. There environment variables goes over configuration files.

* [Examples](#configuration.environmentMapping.examples)
* [Services](#configuration.environmentMapping.services)
* [Global](#configuration.environmentMapping.global)

## Examples {#environmentMapping.examples}

**caretakerd.yaml**
```yaml
services:
    king:
        type: master
        command: ["echo", "Hello world!"]
        user: king
```

**Executions**
```bash
$ caretakerd run
Hello world!

$ export CTD.king.COMMAND=whoami
$ caretakerd run
king

$ export CTD.king.USER=root
$ caretakerd run
root
```

## Services {#environmentMapping.services}

| Environment variable | Configuration property |
| --- | --- |
| ``CTD.<service>.TYPE`` | {@ref github.com/echocat/caretakerd/service.Config#Type} |
| ``CTD.<service>.COMMAND`` | {@ref github.com/echocat/caretakerd/service.Config#Command} |
| ``CTD.<service>.START_DELAY_IN_SECONDS`` | {@ref github.com/echocat/caretakerd/service.Config#StartDelayInSeconds} |
| ``CTD.<service>.RESTART_DELAY_IN_SECONDS`` | {@ref github.com/echocat/caretakerd/service.Config#RestartDelayInSeconds} |
| ``CTD.<service>.SUCCESS_EXIT_CODES`` | {@ref github.com/echocat/caretakerd/service.Config#SuccessExitCodes} |
| ``CTD.<service>.STOP_WAIT_IN_SECONDS`` | {@ref github.com/echocat/caretakerd/service.Config#StopWaitInSeconds} |
| ``CTD.<service>.USER`` | {@ref github.com/echocat/caretakerd/service.Config#User} |
| ``CTD.<service>.DIRECORY`` | {@ref github.com/echocat/caretakerd/service.Config#Directory} |
| ``CTD.<service>.AUTO_RESTART`` | {@ref github.com/echocat/caretakerd/service.Config#AutoRestart} |
| ``CTD.<service>.INHERIT_ENVIRONMENT`` | {@ref github.com/echocat/caretakerd/service.Config#InheritEnvironment} |
| ``CTD.<service>.LOG_LEVEL`` | {@ref github.com/echocat/caretakerd/service.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#Level} |
| ``CTD.<service>.LOG_STDOUT_LEVEL`` | {@ref github.com/echocat/caretakerd/service.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#StdoutLevel} |
| ``CTD.<service>.LOG_STDERR_LEVEL`` | {@ref github.com/echocat/caretakerd/service.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#StderrLevel} |
| ``CTD.<service>.LOG_FILE_NAME`` | {@ref github.com/echocat/caretakerd/service.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#Filename} |
| ``CTD.<service>.LOG_MAX_SIZE_IN_MB`` | {@ref github.com/echocat/caretakerd/service.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#MaxSizeInMb} |
| ``CTD.<service>.LOG_MAX_BACKUPS`` | {@ref github.com/echocat/caretakerd/service.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#MaxBackups} |
| ``CTD.<service>.LOG_MAX_AGE_IN_DAYS`` | {@ref github.com/echocat/caretakerd/service.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#MaxAgeInDays} |
| ``CTD.<service>.ENVIRONMENT.<environmentName>`` | {@ref github.com/echocat/caretakerd/service.Config#Environment}``[<environmentName>]`` |

## Global {#environmentMapping.global}

| Environment variable | Configuration property |
| --- | --- |
| ``CTD.LOG_LEVEL`` | {@ref github.com/echocat/caretakerd.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#Level} |
| ``CTD.LOG_STDOUT_LEVEL`` | {@ref github.com/echocat/caretakerd.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#StdoutLevel} |
| ``CTD.LOG_STDERR_LEVEL`` | {@ref github.com/echocat/caretakerd.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#StderrLevel} |
| ``CTD.LOG_FILE_NAME`` | {@ref github.com/echocat/caretakerd.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#Filename} |
| ``CTD.LOG_MAX_SIZE_IN_MB`` | {@ref github.com/echocat/caretakerd.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#MaxSizeInMb} |
| ``CTD.LOG_MAX_BACKUPS`` | {@ref github.com/echocat/caretakerd.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#MaxBackups} |
| ``CTD.LOG_MAX_AGE_IN_DAYS`` | {@ref github.com/echocat/caretakerd.Config#Logger}: {@ref github.com/echocat/caretakerd/logger.Config#MaxAgeInDays} |
