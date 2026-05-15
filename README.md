# prefwatch

[![Go Report Card](https://goreportcard.com/badge/github.com/donatj/prefwatch)](https://goreportcard.com/report/github.com/donatj/prefwatch)
[![CI](https://github.com/donatj/prefwatch/actions/workflows/ci.yml/badge.svg)](https://github.com/donatj/prefwatch/actions/workflows/ci.yml)

Watch system perference plists for changes and report diffs.

`prefwatch` monitors macOS preference files for changes. It watches `~/Library/Preferences` for plist file modifications. When a preference file changes, it prints a unified diff of the changed values to the console.

## Installation

### Binaries

Signed and notarized binares are available on [releases](https://github.com/donatj/prefwatch/releases).

### Compile

```sh
go install github.com/donatj/prefwatch@latest
```

## Usage

Run the command. It watches `~/Library/Preferences` until terminated.

```sh
prefwatch
```

Output appears in the console each time a preference file changes.
