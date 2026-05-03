# utils

A collection of personal command-line interface tools.

- `logbook` - Manage a simple markdown journal per project
- `ogx` - Inspect OpenGraph tags for a URL
- `forkit` - Embedded implementation prototypes

## Installation

Put this repo’s `bin` directory on your `PATH`.

```plain
cd ~/Developer
git clone git@github.com:kelvinampofo/utils.git
# zsh/bash
export PATH="$PATH:$HOME/Developer/workspaces/utils/bin"

# fish
fish_add_path $HOME/Developer/workspaces/utils/bin
```

## Build

Everything is written in Go. Use the make command to build all utils:

```plain
cd ~/Developer/workspaces/utils
make build-all
```

Build one util:

```plain
make build UTIL=logbook
```

## Usage

Run any tool with `--help` for its commands and flags.
