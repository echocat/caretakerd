# Getting started {#gettingStarted}

* [Installing](#gettingStarted.installing)
* [Configuring](#gettingStarted.configuring)
* [Run](#gettingStarted.run)
* [Control caretakerd with caretakerctl](#gettingStarted.control-caretakerd-with-caretakerctl)

## Installing {#gettingStarted.installing}

Download your matching distribution files and extract it where you want.

Linux 64-Bit example:
```bash
# Download direct from GitHub...
sudo curl -SL https://github.com/echocat/caretakerd/releases/download/v{{.Version}}/caretakerd-linux-amd64.tar.gz \
    | tar -xz --exclude caretakerd.html -C /usr/bin

# Download from caretakerd.echocat.org...   
sudo curl -SLk https://caretakerd.echocat.org/v{{.Version}}/download/caretakerd-linux-amd64.tar.gz \
    | tar -xz --exclude caretakerd.html -C /usr/bin

# Download always the latest from caretakerd.echocat.org...   
sudo curl -SLk https://caretakerd.echocat.org/latest/download/caretakerd-linux-amd64.tar.gz \
    | tar -xz --exclude caretakerd.html -C /usr/bin
```

> Hint: You can use exact these command to in a [Dockerfile](https://docs.docker.com/engine/reference/builder/)
> to create caretakerd in a docker container.

See [Downloads](#downloads) for all possible download locations.

## Configuring {#gettingStarted.configuring}

Place your configuration file at your favorite location. Default location is:

* Linux: ``/etc/caretakerd.yaml``
* Mac OS X: ``/etc/caretakerd.yaml``
* Windows: ``C:\ProgramData\caretakerd\config.yaml``

See [Configuration examples](#configuration.examples) for a quick start or consult [Configuration structure](#configuration.structure) for all possibilities.

## Run {#gettingStarted.run}

```bash
caretakerd run
```

## Control caretakerd with caretakerctl {#gettingStarted.control-caretakerd-with-caretakerctl}

> Precondition: RPC is enabled. See [RPC enabled configuration example](#configuration.examples.rpcEnabled) how to done this.

```bash
# Start a not already running service
$ caretakerctl start peasant

# Retrieve status of a service
$ caretakerctl status peasant

# Stop a not already running service
$ caretakerctl stop peasant
```
