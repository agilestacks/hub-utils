# Hub State API

GCP function backend for [hub-state-api](../api/) to manage states.

## Build

Build is handled by GitHub Action workflow triggered on `push`

## Run/Test locally

```shell
npm install
npm start
```

## Stacks Function filtering examples

### By Stack Status

GET /stacks?status=deployed

### By Initiator

GET /stacks?latestOperation.initiator=akranga

### By Stack Name

GET /stacks?name=GKE

### Deployed Before

GET /stacks?latestOperation.timestamp[before]=2022-05-15

### Deployed After

GET /stacks?latestOperation.timestamp[after]=2022-05-15
