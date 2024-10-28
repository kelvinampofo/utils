# utils

A collection of command-line utilities designed for everyday computer tasks.

- [logbook](bin/logbook) — Manage a simple logbook per project

### Installation

Follow the steps below to set up the utilities on your system:

#### 1. Clone the repository

Clone this repository to your local machine in the ~/Developer directory (or your preferred location):

```bash
mkdir -p ~/Developer
cd ~/Developer
git clone git@github.com:kelvinampofo/utils.git
```

#### 2. Make the scripts executable

Make scripts executable:

```bash
chmod +x ~/Developer/utils/bin/*
```

#### 3. Add the utils to your PATH

Add this line to your shell config file (~/.bashrc or ~/.zshrc):

```bash
export PATH=$PATH:$HOME/Developer/utils/bin
```

Then reload with:

```bash
source ~/.bashrc  # or source ~/.zshrc
```

#### 4. Verify installation

Test the utilities to ensure they are correctly set up. For example, run:

```bash
logbook
```
## Usage

See each program's manual with the `--help` command.