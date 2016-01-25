# ``SocketAddress`` { .property }
Simple value

## Description

This represents a socket address in format ``<protocol>://<target>``.

## Protocols

### ``tcp``

This address connects or binds to a TCP socket. The ``target`` should be of format ``<host>:<port>``.

### ``unix``

This address connects or binds to a UNIX file socket. The ``target`` should be the location of the socket file.
