# ``rpc.Rpc`` { .property }
Structure

## Introduction

The RPC connector handles remote produce calls to caretakerd by [``caretakerctl``](../../executables/caretakerctl) or
the child processes. It makes it able to start, stop, restart, ... every process and other things.

## Properties

### ``enabled`` { #enabled .property }
= ``false`` ([``Boolean``](Boolean))

If this is set to ``true`` caretakerd is accessible by RPC.

### ``listen`` { #listen .property }
= ``tcp://localhost:57955`` ([``SocketAddress``](SocketAddress))

This is the socket where the RPC connector of caretakerd is listen to. This property will be ignored if
[``enabled``](#enabled) is set to ``false``.

### ``securityStore`` { #securityStore .property }
= ([``SecurityStore``](SecurityStore))

Controls how the RPC connector is secured.

*See [``SecurityStore``](SecurityStore) for more details.*
