# Getting started {#gettingStarted}

* [Installing](#gettingStarted.installing)
* [Configuring](#gettingStarted.configuring)
* [Run](#gettingStarted.run)
* [Control caretakerd with caretakerctl](#gettingStarted.control-caretakerd-with-caretakerctl)
* [See in action with Docker](#gettingStarted.see-in-action-with-docker)

## Installing {#gettingStarted.installing}

Download your matching distribution files and extract it where you like.

Linux 64-Bit example:
```bash
# Download directly from GitHub...
sudo curl -SL https://github.com/echocat/caretakerd/releases/download/v{{.Version}}/caretakerd-linux-amd64.tar.gz \
    | tar -xz --exclude caretakerd.html -C /usr/bin

# Download from caretakerd.echocat.org...   
sudo curl -SL https://caretakerd.echocat.org/v{{.Version}}/download/caretakerd-linux-amd64.tar.gz \
    | tar -xz --exclude caretakerd.html -C /usr/bin

# Download always the latest version from caretakerd.echocat.org...   
sudo curl -SL https://caretakerd.echocat.org/latest/download/caretakerd-linux-amd64.tar.gz \
    | tar -xz --exclude caretakerd.html -C /usr/bin
```

> Hint: You can use exactly this command in a [Dockerfile](https://docs.docker.com/engine/reference/builder/)
> to create caretakerd in a docker container.

See [Downloads](#downloads) for all possible download locations.

## Configuring {#gettingStarted.configuring}

Place your configuration file at your favorite location. Default location is:

* Linux: ``/etc/caretakerd.yaml``
* Mac OS X: ``/etc/caretakerd.yaml``
* Windows: ``C:\ProgramData\caretakerd\config.yaml``

See [Configuration examples](#configuration.examples) for a quick start or consult the [Configuration structure](#configuration.structure) for all possibilities.

## Run {#gettingStarted.run}

```bash
caretakerd run
```

## Control caretakerd with caretakerctl {#gettingStarted.control-caretakerd-with-caretakerctl}

> Precondition: RPC is enabled. See [RPC enabled configuration example](#configuration.examples.rpcEnabled) how to done this.

```bash
# Start a service that is not already running 
$ caretakerctl start peasant

# Retrieve the status of a service
$ caretakerctl status peasant

# Stop an already running service
$ caretakerctl stop peasant
```

## See in action with Docker {#gettingStarted.see-in-action-with-docker}

You can find on [github.com/echocat/caretakerd-docker-demos](https://github.com/echocat/caretakerd-docker-demos) lots of working demos. Try it out!
