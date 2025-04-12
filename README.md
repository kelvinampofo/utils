# utils

Minimal CLI utilities for everyday personal computer tasks.

[logbook](cmd/logbook/main.go) — manage a project-specific logbook in markdown.
<br>
[new-sketch](cmd/new-sketch/main.go) — scaffold a new Vite app for rapid UI prototyping.

**Installation**

The idea is to add these utils to your `PATH`, similar to this:

```
cd ~/Developer
git clone git@github.com:kelvinampofo/utils.git

# edit your shell configuration (e.g., ~/.zshrc or ~/.bashrc)
echo 'export PATH=$PATH:$HOME/Developer/utils/bin' >> ~/.zshrc
source ~/.zshrc
```

**Go-based utilities**

The CLI utils are written in Go and must be built before use. If you haven't already, please install [Go](https://go.dev) on your machine. To build them:

```
cd ~/Developer/utils
go build -o bin/<util-name> ./cmd/<util-name>
```

Make sure the resulting `bin/<util-name>` is in your `$PATH` as shown above. You can now run it using:

```
<util-name>
```

## Usage

For a manual on each util, run:

```
<util-name> --help
```
