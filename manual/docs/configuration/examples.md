# Examples

## Simple application

```yaml
# Run the service king and queen at startup.
# queen will run every 5th second until king is finished.
# king will run for 120 seconds. If king is finished the whole careteakerd will
# go down.

services:
    king:
        type: master
        command: ["sleep", "120"]

    queen:
        type: autoStart
        command: ["echo","Hello world!"]
        cronExpression: "0,5,10,15,20,25,30,35,40,45,50,55 * * * * *"
```

## RPC enabled

```yaml
# Run the service king at startup and enable rpc.
# peasant will only be started of "caretakerctl start peasant" is called.
# king will run for 120 seconds. If king is finished the whole careteakerd
# will go down.

rpc:
    enabled: true

services:
    king:
        type: master
        command: ["sleep", "120"]

    peasant:
        type: onDemand
        command: ["sleep.exe","10"]
```
