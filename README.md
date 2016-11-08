[![Circle CI](https://circleci.com/gh/echocat/caretakerd.svg?style=svg)](https://circleci.com/gh/echocat/caretakerd)
[![Go Report Card](https://goreportcard.com/badge/github.com/echocat/caretakerd)](https://goreportcard.com/report/github.com/echocat/caretakerd)
[![Code Climate](https://codeclimate.com/github/echocat/caretakerd/badges/gpa.svg)](https://codeclimate.com/github/echocat/caretakerd)
[![Gitter](https://badges.gitter.im/echocat/caretakerd.svg)](https://gitter.im/echocat/caretakerd?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

# caretakerd

caretakerd is a simple process supervisor. There are no external dependencies and it is optimized for containerization (like Docker) and simple configuration.

* [Documentation](#documentation)
* [Building](#building)
* [Contributing](#contributing)
* [Support](#support)
* [License](#license)

## Documentation

For general documentation, please refer to [caretakerd.echocat.org](https://caretakerd.echocat.org).

For specific versions, please refer to [caretakerd.echocat.org/all](https://caretakerd.echocat.org/all).

## Building

### Requirements

To build caretakerd, you only need:

* a compatible operating system (Linux, Windows or Mac OS X)
* a working Java 8 installation

The build system will download every dependency and build it if necessary.

> **Hint:** The Go runtime build by the build system will be stored under ``~/.go``.

### Run

To run caretakerd on Linux and Mac OS X, invoke the following:

```bash
# Build binaries only
./gradlew build

# Run tests (includes compile)
./gradlew test

# Build resulting packages (includes documentation - includes compile)
./gradlew package

# Set the target version number, increase the version number, do gradlew package,
# deploy everything to GitHub releases and set next development version number.
./gradlew release
```

To run caretakerd on Windows, invoke the following:

```bash
# Build binaries only
gradlew build

# Run tests (includes compile)
gradlew test

# Build resulting packages (includes documentation - includes compile)
gradlew package

# Set the target version number, increase the version number, do gradlew package,
# deploy everything to GitHub releases and set next development version number.
gradlew release
```

### Build artifacts

* You can find the compiled and linked binaries under ``./build/out/caretaker*``
* You can find the generated document under ``./build/docs/caretakerd.html`` and ``./build/docs/caretakerd-windows.html`` 
* You can find the packaged TARZs and ZIPs under ``./build/distribution/caretakerd-*.tar.gz`` and ``./build/distribution/caretakerd-*.zip``

## Contributing

caretakerd is an open source project by [echocat](https://echocat.org).
So if you want to make this project even better, you can contribute to this project on [Github](https://github.com/echocat/caretakerd)
by [fork us](https://github.com/echocat/caretakerd/fork).

If you commit code to this project, you have to accept that this code will be released under the [license](#license) of this project.

## Support

If you need support you can create a ticket in our [issue tracker](https://github.com/echocat/caretakerd/issues)
or join our chat at [echocat.slack.com/messages/caretakerd](https://echocat.slack.com/messages/caretakerd/).

## License

See the [LICENSE](LICENSE) file.
