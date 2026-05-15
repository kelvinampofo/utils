package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	logDir       = ".logbook"
	logExtension = ".md"
	dateFmt      = "2006-01-02"
)

const help = `logbook keeps project-local notes in ./.logbook/YYYY-MM-DD.md.

Usage:
  logbook             open today's log in $EDITOR
  logbook add <text>  append text to today's log
  logbook read [date] print a log
  logbook ls          list logs
  logbook path [date] print a log path
  logbook grep <term> search logs
  logbook help        show this help

Dates use the ISO 8601 calendar format: YYYY-MM-DD.`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "logbook: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return edit("")
	}

	switch args[0] {
	case "-h", "--help", "help":
		fmt.Println(help)
	case "edit":
		date, err := optionalDate(args[1:])
		if err != nil {
			return err
		}
		return edit(date)
	case "add":
		return add(args[1:])
	case "read":
		date, err := optionalDate(args[1:])
		if err != nil {
			return err
		}
		return read(date)
	case "ls", "list":
		return list()
	case "path":
		date, err := optionalDate(args[1:])
		if err != nil {
			return err
		}
		file, err := logFile(date)
		if err != nil {
			return err
		}
		fmt.Println(file)
	case "grep":
		return grep(args[1:])
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], help)
	}
	return nil
}

func add(words []string) error {
	if len(words) == 0 {
		return edit("")
	}
	dir, err := logDirPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	file, err := logFile("")
	if err != nil {
		return err
	}
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, strings.Join(words, " "))
	return err
}

func edit(date string) error {
	dir, err := logDirPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	file, err := logFile(date)
	if err != nil {
		return err
	}
	if err := touch(file); err != nil {
		return err
	}
	cmd, err := editorCommand(file)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func read(date string) error {
	file, err := logFile(date)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	if stdoutIsTTY() && bytes.Count(content, []byte{'\n'}) > 20 {
		return command(defaultPager(), file).Run()
	}
	_, err = os.Stdout.Write(content)
	return err
}

func list() error {
	dir, err := logDirPath()
	if err != nil {
		return err
	}
	files, err := filepath.Glob(filepath.Join(dir, "*"+logExtension))
	if err != nil {
		return err
	}
	for _, file := range files {
		fmt.Println(strings.TrimSuffix(filepath.Base(file), logExtension))
	}
	return nil
}

func grep(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("grep needs a search term")
	}
	dir, err := logDirPath()
	if err != nil {
		return err
	}
	return command("grep", append([]string{"-Rin"}, append(args, dir)...)...).Run()
}

func optionalDate(args []string) (string, error) {
	if len(args) == 0 {
		return "", nil
	}
	if len(args) > 1 {
		return "", fmt.Errorf("too many arguments")
	}
	if _, err := time.Parse(dateFmt, args[0]); err != nil {
		return "", fmt.Errorf("invalid date %q; expected ISO 8601 format YYYY-MM-DD", args[0])
	}
	return args[0], nil
}

func logDirPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	return filepath.Join(wd, logDir), nil
}

func logFile(date string) (string, error) {
	if date == "" {
		date = time.Now().Format(dateFmt)
	}
	dir, err := logDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, date+logExtension), nil
}

func touch(file string) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	return f.Close()
}

func command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func editorCommand(file string) (*exec.Cmd, error) {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return command(editor, file), nil
	}

	for _, editor := range []string{"nvim", "vim"} {
		if _, err := exec.LookPath(editor); err == nil {
			return command(editor, file), nil
		}
	}

	if _, err := exec.LookPath("open"); err == nil {
		textEdit := command("open", "-e", file)
		return textEdit, nil
	}

	return nil, fmt.Errorf("set EDITOR or install nvim or vim")
}

func defaultPager() string {
	if pager := os.Getenv("PAGER"); pager != "" {
		return pager
	}
	return "less"
}

func stdoutIsTTY() bool {
	info, err := os.Stdout.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
