# Features {#features}

* **[Simple configuration](#configuration.examples)**<br>
  Configure everything via one [YAML](https://en.wikipedia.org/wiki/YAML) file. Every configuration could also be overwritten via
  [environment variables](#configuration.environmentMapping) and the master service could
  [receive its command line also from the command line of caretakerd itself](#configuration.dataType.service.Service.command.special-master-handling).
  
* **Optimized for containerization**<br>
  cratekerd was designed to be a simple and small process supervisor for containerization environment (such as [Docker](https://en.wikipedia.org/wiki/Docker_\(software\))).
  It is explicitly not a replacement for host process supervisors like [systemd](https://en.wikipedia.org/wiki/Systemd).
  
* **No dependencies**<br>
  Just [download the binary](#downloads) and use it. caretakerd is a [fat binary](https://en.wikipedia.org/wiki/Fat_binary) ready to use.
  You are not required to install environment with a lot of dependencies like Ruby, Python or Perl - caretakerd is build with [Go](https://en.wikipedia.org/wiki/Go_\(programming_language\)).
    
* **Builtin watchdog**<br>
  If a service crashes, caretakerd will [restart](#configuration.dataType.service.Service.autoRestart) it for you.
  
* **[Custom logging](configuration.dataType.logger.Logger)**<br>
  You can log how you want. It is possible to log to [log files](#configuration.dataType.logger.Logger.filename)
  per service, to just one log file for all services or also direct to [console](#configuration.dataType.logger.Logger.filename).
  Also log rotate based on [max file size](#configuration.dataType.logger.Logger.maxSizeInMb) and [age of](#configuration.dataType.logger.Logger.maxAgeInDays)
  log files are builtin.
  
* **[Builtin cron](#configuration.dataType.service.CronExpression)**<br>
  Execute services when and how often you want.
  
* **[Focus on one core service](#configuration.dataType.service.Type)**<br>
  There is always one service (the [master](#configuration.dataType.service.Type.master))
  which all other services are life and die together. We expect that this service is the most important one and all other services have to serve him.
  
* **[Remote controllable](#configuration.dataType.rpc.Rpc)**<br>
  caretakerd could fully controlled via [``caretakerctl``](#commands.caretakerctl) command or
  REST calls via encrypted HTTPS.