# utils

A lightweight collection of command-line utilities designed to simplify everyday computer tasks.

[logbook](bin/logbook) — Manage a project-specific logbook in markdown.

**Installation**

The idea is to add these utils to your `PATH`, similar to this:

```
cd ~/Developer
git clone git@github.com:kelvinampofo/utils.git
# edit your shell configuration (e.g., ~/.zshrc or ~/.bashrc)
echo 'export PATH=$PATH:$HOME/Developer/utils/bin' >> ~/.zshrc
source ~/.zshrc
```

## Usage

For a manual on each util, run:

```
<util-name> --help
```

Some notes:

### logbook

`logbook` can be configured using an optional `.logbookrc` file, which is recommended for specifying defaults such as the project directory and editor. If this file is not present, the program will revert to sane defaults or command-line arguments respectively.

Copy the example and customise it:

```
cp .logbookrc.example .logbookrc

# replace the `cp` command with `copy` in command prompt or `Copy-Item` in powershell, respectively
```

You can configure the following options:

- `PROJECT_DIR`:
  Path to the default project directory where the logbook will store and manage logs. If not specified, the current directory is used by default.

- `EDITOR`:
  Default text editor to use for editing log files. The utility will fall back to vim or nano if this option is not set and the specified editor is unavailable.
