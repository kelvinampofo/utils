# utils

A toolbox of small CLI programs for everyday tasks.

## Installation

Build the tools and put this repo's `bin` directory on your `PATH`.

```plain
cd ~/Developer/workspaces/utils
make build

# zsh/bash:
export PATH="$PATH:$HOME/Developer/workspaces/utils/bin"

# fish:
fish_add_path $HOME/Developer/workspaces/utils/bin
```

## Build

```plain
make list   # list available tools
make build  # build Go tools into ./bin
make clean  # remove ./bin
```
