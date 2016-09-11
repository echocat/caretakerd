# Features {#features}

* **[Simple configuration](#configuration.examples)**<br>
  Configure everything via one [YAML](https://en.wikipedia.org/wiki/YAML) file. Every configuration could also be overwritten via
  [environment variables](#configuration.environmentMapping) and the master service could
  [also receive its command line from the command line of caretakerd itself](#configuration.dataType.service.Service.command.special-master-handling).
  
* **Optimized for containerization**<br>
  caretakerd was designed to be a simple and small process supervisor for containerization environments (such as [Docker](https://en.wikipedia.org/wiki/Docker_\(software\))).
  It is no replacement for host process supervisors like [systemd](https://en.wikipedia.org/wiki/Systemd).
  
* **No dependencies**<br>
  Just [download the binary](#downloads) and use it. caretakerd is a [fat binary](https://en.wikipedia.org/wiki/Fat_binary) ready to use.
  You are not required to install an environment with a lot of dependencies like Ruby, Python or Perl - caretakerd is build with [Go](https://en.wikipedia.org/wiki/Go_\(programming_language\)).
    
* **Builtin watchdog**<br>
  If a service crashes, caretakerd will [restart](#configuration.dataType.service.Service.autoRestart) it for you.
  
* **[Custom logging](#configuration.dataType.logger.Logger)**<br>
  You can log to your liking. It is possible to [log to files](#configuration.dataType.logger.Logger.filename)
  per service, to just log one file for all services or to directly log to the [console](#configuration.dataType.logger.Logger.filename).
  Also log rotation based on [max file size](#configuration.dataType.logger.Logger.maxSizeInMb) and [age of](#configuration.dataType.logger.Logger.maxAgeInDays)
  log files are built in.
  
* **[Builtin cron](#configuration.dataType.service.CronExpression)**<br>
  Execute services when and how often you want.
  
* **[Focus on one core service](#configuration.dataType.service.Type)**<br>
  There is always one ([master service](#configuration.dataType.service.Type.master)).
  All other services life and die together. We expect that this service is the most important one and all other services have to serve the master.
  
* **[Remote controllable](#configuration.dataType.rpc.Rpc)**<br>
  caretakerd can be fully controlled via [``caretakerctl``](#commands.caretakerctl) command or
  REST calls via encrypted HTTPS.
