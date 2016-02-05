# caretakerd

caretakerd is a minimal process supervisor build for easy use with no dependencies.

* [User documentation](#user-documentation)
* [Building](#building)
* [Contributing](#contributing)
* [Support](#support)
* [License](#license)

## User documentation

Each release of caretakerd is shipping its own documentation.
Please go to the [release page](https://github.com/echocat/caretakerd/releases/latest)
and download the ``caretakerd.html`` for the full documentaion or take the platform specific one out
of the release package (``caretakerd-*.tar.gz`` or ``caretakerd-*.zip``) on the [release page](https://github.com/echocat/caretakerd/releases/latest).

## Building

### Precondition

For building caretakerd there is only a compatible operating system (Linux, Windows or Mac OS X) and a Java 8 installation required.

### Run

On Linux and Mac OS X:
```bash
# Build binaries only
./mvnw compile

# Run tests (includes compile)
./mvnw test

# Build resulting packages (including documentation - includes compile)
./mvnw package

# Set the target version number, increase the version number, do mvnw package,
# deploy everything to GitHub releases and set next development version number.
./mvnw release:prepare release:perform
```

On Windows:
```bash
# Build binaries only
mvnw compile

# Run tests (includes compile)
mvnw test

# Build resulting packages (including documentation - includes compile)
mvnw package

# Set the target version number, increase the version number, do mvnw package,
# deploy everything to GitHub releases and set next development version number.
mvnw release:prepare release:perform
```

### Build artifacts

* Compiled and lined binaries can be found under ``./target/gopath/bin/caretaker*``
* Generated document can be found under ``./target/docs/caretakerd.html``
* Packaged TARZs and ZIPs can be found under ``./target/caretakerd-*.tar.gz`` and ``./target/caretakerd-*.zip``

## Contributing

caretakerd is an open source project of [echocat](https://echocat.org).
So if you want to make this project even better, you can contribute to this project on [Github](https://github.com/echocat/caretakerd)
by [fork us](https://github.com/echocat/caretakerd/fork).

If you commit code to this project you have to accept that this code will be released under the [license](#license) of this project.

## Support

If you need support you can file a ticket at our [issue tracker](https://github.com/echocat/caretakerd/issues)
or join our chat at [echocat.slack.com/messages/caretakerd](https://echocat.slack.com/messages/caretakerd/).

## License

See [LICENSE](LICENSE) file.
