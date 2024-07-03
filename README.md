# Jetscript

Jetscript is a fairly simplistic yet powerful mapping language to implement Single Message Transformations (SMTs).
It is designed to be embedded into NATS, but can be used in any context where a simple mapping language is needed.

## Features
- Simple syntax
- Support for maps; the equivalent of a function in other languages
- One incoming message to 0-many resulting messages 

## CLI Features
- *exec* - Subscribe to a nats subject and execute a jetscript file on each message
- *lint* - Lint a jetscript file for syntax errors
- *push* - Push a jetscript file to a nats the Nats Jetstream Object Store
- *pull* - Pull a jetscript file from the Nats Jetstream Object Store

## Usage
You can write your jetscript in a file locally and lint it with the `jetstream lint <your-file>` command. If you want
to apply your jetscript to messages in a nats subject, you first need to push it to the Jetstream Object Store with
the `jetstream push --path <path-in-object-store> <your-file>` command. Then launch the 
`jetstream exec --subject <subject> --path <path-in-object-store>` command to start processing messages.