# prefwatch

[![CI](https://github.com/donatj/prefwatch/actions/workflows/ci.yml/badge.svg)](https://github.com/donatj/prefwatch/actions/workflows/ci.yml)

prefwatch monitors macOS preference files for changes. It watches `~/Library/Preferences` for plist file modifications. When a preference file changes, it prints a unified diff of the changed values to the console.

## Installation

```
go install github.com/donatj/prefwatch@latest
```

## Usage

Run the command. It watches `~/Library/Preferences` until terminated.

```
prefwatch
```

Output appears in the console each time a preference file changes.
