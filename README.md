# utils

A lightweight collection of command-line utilities designed to simplify everyday computer tasks.

- [logbook](bin/logbook) — manage a project-specific logbook in markdown.
- [cmpress](bin/cmpress) — simplify video and image compression tasks using FFmpeg.
- [new-sandbox](bin/new-sandbox) — quickly create a new Vite app for rapid prototyping.

**Installation**

The idea is to add these utils to your `PATH`, similar to this:

```
cd ~/Developer
git clone git@github.com:kelvinampofo/utils.git
# edit your shell configuration (e.g., ~/.zshrc or ~/.bashrc)
echo 'export PATH=$PATH:$HOME/Developer/utils/bin' >> ~/.zshrc
source ~/.zshrc
```

**Building Go-based utilities**

Some tools, like `logbook`, are written in Go and must be built before use. To build them:

```
cd ~/Developer/utils
go build -o bin/logbook ./cmd/logbook
```

Make sure the resulting `bin/logbook` is in your `$PATH` as shown above. You can now run it using:

```
logbook --help
```

## Usage

For a manual on each util, run:

```
<util-name> --help
```

Some notes:

### cmpress

Ensure `FFmpeg` is installed on your system to use `cmpress`. Follow [FFmpeg installation instructions](https://ffmpeg.org/download.html) if needed.
