# utils

A collection of command-line utilities designed to streamline everyday computer tasks.

- [logbook](bin/logbook) — Manage a simple logbook per project

## Installation

Follow the steps below to set up the utilities on your system:

### 1. Clone the repository

Clone this repository to your local machine in the ~/Developer directory (or your preferred location):

```bash
mkdir -p ~/Developer
cd ~/Developer
git clone https://github.com/kelvinampofo/utils.git
```

### 2. Make the scripts executable

Ensure all scripts in the bin directory are executable:

```bash
chmod +x ~/Developer/utils/bin/*
```

### 3. Add the utils to your PATH

To access the utilities from anywhere, add the bin directory to your PATH:
  
1. Open your shell configuration file (~/.bashrc for Bash, ~/.zshrc for Zsh):
 
  ```bash
  vim ~/.bashrc
  ```

  or for Zsh:

  ```bash
  vim ~/.zshrc
  ```
2. Add the following line to the file:
  
  ```bash
export PATH=$PATH:$HOME/Developer/utils/bin
  ```
3. Save and exit the editor.
  
4. Reload your shell configuration:
  
  ```bash
source ~/.bashrc
  ```

  or for Zsh:

  ```bash
source ~/.zshrc
  ```

4. Verify the installation

Test the utilities to ensure they are correctly set up. For example, run:

```bash
logbook
```

If the installation is successful, this command will create a .logbook directory in your current working directory and open today’s log file with your default editor (e.g. vim).

## Usage

See each program's manual with `--help`