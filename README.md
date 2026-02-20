# utils

CLI utils for everyday computer things.

- `logbook` — Manage a simple markdown journal per project
- `ogx` — Check OpenGraph tags for a URL

## Installation

Add this repo's `bin` directory to your `PATH`.

```plain
cd ~/Developer
git clone git@github.com:kelvinampofo/utils.git

# zsh/bash
export PATH="$PATH:$HOME/Developer/workspaces/utils/bin"

# fish
set -U fish_user_paths $HOME/Developer/workspaces/utils/bin $fish_user_paths
```

## Build

Utils are written in Go.

```plain
cd ~/Developer/workspaces/utils
make build-all
```

Build one util:

```plain
make build UTIL=logbook
```

## Usage

See each utils help with `<util-name> --help`.
