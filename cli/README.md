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

## Release process

Release process is implemented via [GoReleaser] and GitHub Actions. [GoReleaser] is responsible for building, packing and creation of GitHub releases. GitHub Actions triggers this process when tag with name in [semver](https://semver.org/) format is created. So to trigger Packaging process, for example, you need to run next command:

```bash
git tag -a v1.0.0 -m "First release"
git push origin v1.0.0
```

[GoReleaser]: https://goreleaser.com/
