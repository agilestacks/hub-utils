# Hub State extensions

[Hub CLI](https://github.com/agilestacks/hub) extension which work together with [hub-state-api](../hub-state-api/) to manage states.

## Usage example

```shell
hub state ls
hub state ls --filter="latestOperation.timestamp[after]=2022-05-19,status=incomplete,initiator=akranga"
hub state ls --filter="latestOperation.timestamp[before]=2022-05-19"
hub state show <sandbox id>
```

## Build

```shell
go build -o bin/hub-state main.go
```

## Run

```shell
go run main.go
```

