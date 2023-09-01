# Command: `hub-state rm`

Removes stack state from the project by Stack ID

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
hub-state rm <Stack ID>
```

## See also

* [`hub-state show`](hub-state-show.md)
* [`hub-state ls`](hub-state-ls.md)
