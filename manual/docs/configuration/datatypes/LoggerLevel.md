# ``LoggerLevel`` { .property }
Enum

## Description

Represents a level for logging with a [``Logger``](#Logger).

## Values

### ``debug`` { #debug }

Used for debugging proposes. This level is only required you something goes wrong and you need more information.

### ``info`` { #info }

This is the regular level. Every normal message will be logged with this level.

### ``warning`` { #warning }

If a problem appears but the program is still able to continue its work, this level is used.

### ``error`` { #error }

If a problem appears and the program is not longer able to continue its work, this level is used.

### ``fatal`` { #fatal }

This level is used on dramatic problems.
