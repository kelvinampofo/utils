package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Config struct {
	LogDir          string
	FileExtension   string
	FilePermission  os.FileMode
	EntryPermission os.FileMode
}

var config = Config{
	LogDir:          ".logbook",
	FileExtension:   ".md",
	FilePermission:  0755,
	EntryPermission: 0644,
}

// fatal prints an error message and exits the program.
func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

// runCmd sets up the command's I/O and executes it. If execution fails, it prints an error message and exits.
func runCmd(cmd *exec.Cmd, msg string) {
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fatal("%s: %v", msg, err)
	}
}

func main() {
	rootCmd := newRootCmd()
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fatal("%v", err)
	}
}

// newRootCmd initialises the root command and registers all subcommands.
func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "logbook [subcommand]",
		Short: "A CLI for daily Markdown logbooks.",
		Run:   func(cmd *cobra.Command, args []string) { cmd.Help() },
	}
	rootCmd.AddCommand(addCmd, editCmd, readCmd, lsCmd, grepCmd, logfileCmd, logdirCmd)
	return rootCmd
}

var editCmd = &cobra.Command{
	Use:   "edit [date]",
	Short: "Open today's log entry or a specified date file",
	Run: func(cmd *cobra.Command, args []string) {
		editEntry(resolveLogFile(firstArg(args)))
	},
}

var addCmd = &cobra.Command{
	Use:   "add [text]",
	Short: "Append a line of text to today's log",
	Run: func(cmd *cobra.Command, args []string) {
		file := resolveLogFile("")
		if len(args) == 0 {
			editEntry(file)
		} else {
			appendToEntry(file, args)
		}
	},
}

var readCmd = &cobra.Command{
	Use:   "read [date]",
	Short: "Read today's or a specified date log in a pager",
	Run: func(cmd *cobra.Command, args []string) {
		readEntry(resolveLogFile(firstArg(args)))
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all log files in the .logbook directory",
	Run: func(cmd *cobra.Command, args []string) {
		listEntries(resolveLogDir())
	},
}

var grepCmd = &cobra.Command{
	Use:   "grep <keyword>",
	Short: "Search logs for matching lines",
	Run: func(cmd *cobra.Command, args []string) {
		grepEntries(resolveLogDir(), args)
	},
}

var logfileCmd = &cobra.Command{
	Use:   "logfile",
	Short: "Print path to today's log file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(resolveLogFile(""))
	},
}

var logdirCmd = &cobra.Command{
	Use:   "logdir",
	Short: "Print path to the .logbook directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(resolveLogDir())
	},
}

// firstArg returns the validated date from args or today's date in YYYY-MM-DD format.
func firstArg(args []string) string {
	if len(args) > 0 {
		if _, err := time.Parse("2006-01-02", args[0]); err != nil {
			fatal("invalid date format: %s (expected YYYY-MM-DD)", args[0])
		}
		return args[0]
	}
	return time.Now().Format("2006-01-02")
}

// resolveLogDir returns the full path to the .logbook directory in the current working directory.
func resolveLogDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fatal("error getting current dir: %v", err)
	}
	return filepath.Join(dir, config.LogDir)
}

// resolveLogFile returns the full path to the log file for the given date.
func resolveLogFile(date string) string {
	return filepath.Join(resolveLogDir(), date+config.FileExtension)
}

// editEntry opens the specified log file in the user's default editor.
func editEntry(file string) {
	ensureDir(filepath.Dir(file), config.FilePermission)
	runCmd(exec.Command(defaultEditor(), file), "failed to open editor")
}

// appendToEntry appends the provided text to the specified log file.
func appendToEntry(file string, lines []string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, config.EntryPermission)
	if err != nil {
		fatal("failed to open file: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(strings.Join(lines, " ") + "\n"); err != nil {
		fatal("failed to write to file: %v", err)
	}
	fmt.Fprintf(os.Stdout, "added entry to %q\n", filepath.Base(file))
}

// readEntry displays the log file using the user's pager, defaulting to "less" if not specified.
func readEntry(file string) {
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "less"
	}
	runCmd(exec.Command(pager, file), "failed to open pager")
}

// listEntries prints the names of all log files in the log directory.
func listEntries(dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*"+config.FileExtension))
	if err != nil {
		fatal("error listing files: %v", err)
	}
	for _, f := range files {
		fmt.Fprintln(os.Stdout, filepath.Base(f))
	}
}

// grepEntries performs a recursive, case-insensitive search in the log directory using grep.
func grepEntries(dir string, args []string) {
	if len(args) == 0 {
		fatal("nothing to grep.")
	}
	grepArgs := append([]string{"-iR", "--color"}, append(args, dir)...)
	runCmd(exec.Command("grep", grepArgs...), "grep failed")
}

// defaultEditor returns the editor set by the EDITOR environment variable or the first available editor from a list.
func defaultEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	for _, editor := range []string{"nvim", "vim", "nano", "emacs"} {
		if _, err := exec.LookPath(editor); err == nil {
			return editor
		}
	}
	fatal("no suitable editor found. Set the EDITOR environment variable.")
	return ""
}

// ensureDir checks if a directory exists and creates it with the specified permissions if not.
func ensureDir(dir string, perm os.FileMode) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, perm); err != nil {
			fatal("could not create directory %s: %v", dir, err)
		}
	}
}
