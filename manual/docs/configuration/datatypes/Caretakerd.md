# ``Caretakerd`` { .property }
Structure

## Description

Central configuration object of caretakerd.

## Properties

### ``rpc`` { #rpc .property }
([``Rpc``](Rpc))

Settings how to access caretakerd via RPC.

*See [``Rpc``](Rpc) for more details.*

### ``control`` { #control .property }
([``Control``](Control))

Access settings for [``caretakerctl``](../../executables/caretakerctl) in caretakerd.

*See [``Control``](Control) for more details.*

### ``logger`` { #logger .property }
([``Logger``](Logger))

Logging settings for logging of the caretakerd itself - this does not affect logging of services.

*See [``Logger``](Logger) for more details.*

### ``services`` { #services .property }
([``[]Service``](Service))

List of all services to run by caretakerd.

!!! important "This list should contain exact one service of type [``master``](ServiceType#master)."

*See [``Service``](Service) for more details.*
