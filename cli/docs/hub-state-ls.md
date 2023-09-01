# Command: `hub-state ls`

List all stacks states within a project

## Flags

| Flag      | Description | Required |
| :-------- | :--------   | :-:      |
| `--filter strings` | Filter by name, status or initiator. Example: --filter "name=GKE,status=incomplete" | |

## Global Flags

These parameters applies across all extension commands

| Flag      | Description | Required |
| :-------- | :--------   | :-:      |
| `-h, --help` | Print help and usage message | |
| `-o, --output string` | Output format. Must be one of [table, json] | |
| `-p, --project string` | GCP Project ID | |
| `-l, --stateAPILocation string` | Location of State API endpoint (default "us-central1") | |
| `-v, --verbose` | Verbose output for diagnostics | |

## Usage examples

```bash
hub-state ls
hub-state ls --filter="latestOperation.timestamp[after]=2022-05-19,status=incomplete,initiator=akranga"
hub-state ls --filter="latestOperation.timestamp[before]=2022-05-19"
```

## See also

* [`hub-state show`](hub-state-show.md)
* [`hub-state rm`](hub-state-rm.md)
