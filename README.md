[![Continuous Integration](https://github.com/echocat/caretakerd/actions/workflows/ci.yml/badge.svg?event=push)](https://github.com/echocat/caretakerd/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/echocat/caretakerd)](https://goreportcard.com/report/github.com/echocat/caretakerd)
[![Code Climate](https://codeclimate.com/github/echocat/caretakerd/badges/gpa.svg)](https://codeclimate.com/github/echocat/caretakerd)

# caretakerd

caretakerd is a simple process supervisor. There are no external dependencies and it is optimized for containerization (like Docker) and simple configuration.

* [Documentation](#documentation)
* [Building](#building)
* [Contributing](#contributing)
* [Support](#support)
* [License](#license)

## Documentation

For general documentation, please refer to [caretakerd.echocat.org](https://caretakerd.echocat.org).

For specific versions, please refer to [caretakerd.echocat.org/all/](https://caretakerd.echocat.org/all/).

## Building

### Requirements

To build caretakerd, you only need:

* a compatible operating system (Linux, Windows or Mac OS X)
* a working Go (at least version 1.11)

The build system will download every dependency and build it if necessary.

### Run

To run caretakerd, invoke the following:

```bash
# Run all tests
go run ./build test

# Build binaries
go run ./build build
```

### Build artifacts

* You can find the compiled and linked binaries under `./var/binaries/`
* You can find the generated document under `./var/manuals/`
* You can find the packaged TARZs and ZIPs under `./var/dist/`

## Contributing

caretakerd is an open source project by [echocat](https://echocat.org).
So if you want to make this project even better, you can contribute to this project on [GitHub](https://github.com/echocat/caretakerd)
by [fork us](https://github.com/echocat/caretakerd/fork).

If you commit code to this project, you have to accept that this code will be released under the [license](#license) of this project.

## Support

If you need support you can create a ticket in our [issue tracker](https://github.com/echocat/caretakerd/issues)
or join our chat at [echocat.slack.com/messages/caretakerd](https://echocat.slack.com/messages/caretakerd/).

## License

See the [LICENSE](LICENSE) file.
