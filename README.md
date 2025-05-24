# utils

Minimal CLI utilities for everyday computer tasks.

- `logbook` — maintain a markdown-based project logbook

**Installation**

Clone the repository and add it to your `PATH`:

```
cd ~/Developer
git clone git@github.com:kelvinampofo/utils.git

# add to PATH (e.g., in ~/.zshrc or ~/.bashrc)
echo 'export PATH="$HOME/Developer/utils/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

**Building Utilities**

Utilities are written in Go and must be built before use. First, install [Go](https://go.dev) if you haven’t already.

Then build a util like this:

```
cd ~/Developer/utils
go build -o bin/<util-name> ./cmd/<util-name>
```

Make sure the `bin/` directory is in your `PATH`. Then run any util with:

```
<util-name>
```

## Usage

For a manual on each util, run:

```
<util-name> --help
```
