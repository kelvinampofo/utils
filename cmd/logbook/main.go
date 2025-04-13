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
	FilePermission:  0755, // rwxr-xr-x
	EntryPermission: 0644, // rw-r--r--
}

func main() {
	rootCmd := newRootCmd()
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// newRootCmd initialises the root command and registers all subcommands.
func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "logbook [subcommand]",
		Short: "A command-line interface for daily Markdown logbooks.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(
		addCmd,
		editCmd,
		readCmd,
		lsCmd,
		grepCmd,
		logfileCmd,
		logdirCmd,
	)

	return rootCmd
}

// editCmd opens today's log entry in the user's default editor or a specified date file.
var editCmd = &cobra.Command{
	Use:   "edit [date]",
	Short: "Open today's log entry (default) or a specified date file",
	Run: func(cmd *cobra.Command, args []string) {
		date := ""
		if len(args) > 0 {
			date = args[0]
		}
		logFile := resolveLogFile(date)
		editEntry(logFile)
	},
}

// addCmd appends a line of text to today's log entry or opens the editor if no text is provided.
var addCmd = &cobra.Command{
	Use:   "add [text]",
	Short: "Append a line of text to today's log",
	Run: func(cmd *cobra.Command, args []string) {
		logFile := resolveLogFile("")
		if len(args) == 0 {
			editEntry(logFile)
			return
		}
		appendToEntry(logFile, args)
	},
}

// readCmd shows today's log in a pager or a specified file
var readCmd = &cobra.Command{
	Use:   "read [date]",
	Short: "Read today's log (default) or a specified date file in a pager",
	Run: func(cmd *cobra.Command, args []string) {
		date := ""
		if len(args) > 0 {
			date = args[0]
		}
		logFile := resolveLogFile(date)
		readEntry(logFile)
	},
}

// lsCmd lists all .md files in the log directory
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all log files in the .logbook directory",
	Run: func(cmd *cobra.Command, args []string) {
		logDir := resolveLogDir()
		listEntries(logDir)
	},
}

// grepCmd searches for strings in the log directory
var grepCmd = &cobra.Command{
	Use:   "grep <keyword>",
	Short: "Search logs for matching lines",
	Run: func(cmd *cobra.Command, args []string) {
		logDir := resolveLogDir()
		grepEntries(logDir, args)
	},
}

// logfileCmd prints the path to today's log file
var logfileCmd = &cobra.Command{
	Use:   "logfile",
	Short: "Print path to today's log file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(resolveLogFile(""))
	},
}

// logdirCmd prints the path to the log directory
var logdirCmd = &cobra.Command{
	Use:   "logdir",
	Short: "Print path to the .logbook directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(resolveLogDir())
	},
}

// resolveLogDir returns the full path to the .logbook directory in the current working directory.
func resolveLogDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting current dir: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(currentDir, config.LogDir)
}

// resolveLogFile returns the full path to the log file for a given date.
// If dateString is empty, it uses today's date in the format YYYY-MM-DD.
func resolveLogFile(dateString string) string {
	if dateString == "" {
		dateString = time.Now().Format("2006-01-02")
	} else {
		if _, err := time.Parse("2006-01-02", dateString); err != nil {
			fmt.Fprintf(os.Stderr, "invalid date format: %s (expected YYYY-MM-DD)\n", dateString)
			os.Exit(1)
		}
	}
	return filepath.Join(resolveLogDir(), dateString+config.FileExtension)
}

// editEntry opens a log file in the user's default editor.
func editEntry(file string) {
	dir := filepath.Dir(file)
	ensureDir(dir, config.FilePermission)

	editor := defaultEditor()
	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to open editor: %v\n", err)
		os.Exit(1)
	}
}

// appendToEntry appends the provided text to the specified log file.
func appendToEntry(file string, lines []string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, config.EntryPermission)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := f.WriteString(strings.Join(lines, " ") + "\n"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "added entry to \"%s\"\n", filepath.Base(file))
}

// readEntry displays the log file using the user's pager, defaulting to 'less' if not specified.
func readEntry(file string) {
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "less"
	}

	cmd := exec.Command(pager, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to open pager: %v\n", err)
		os.Exit(1)
	}
}

// listEntries prints the names of all log files in the log directory.
func listEntries(dir string) {
	files, err := filepath.Glob(filepath.Join(dir, "*"+config.FileExtension))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing files: %v\n", err)
		os.Exit(1)
	}

	for _, f := range files {
		fmt.Fprintln(os.Stdout, filepath.Base(f))
	}
}

// grepEntries performs a recursive case-insensitive search in the log directory using grep.
func grepEntries(dir string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "nothing to grep.")
		os.Exit(1)
	}

	grepArgs := append([]string{"-iR", "--color"}, append(args, dir)...)
	cmd := exec.Command("grep", grepArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "grep failed: %v\n", err)
		os.Exit(1)
	}
}

// defaultEditor returns the editor set by the EDITOR environment variable or the first available editor from a predefined list.
func defaultEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	editors := []string{"nvim", "vim", "nano", "emacs"}
	for _, editor := range editors {
		if _, err := exec.LookPath(editor); err == nil {
			return editor
		}
	}

	fmt.Fprintln(os.Stderr, "no suitable editor found. Set the EDITOR environment variable.")
	os.Exit(1)
	return ""
}

// ensureDir checks if a directory exists and creates it with the specified permissions if it does not.
func ensureDir(dir string, perm os.FileMode) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, perm); err != nil {
			fmt.Fprintf(os.Stderr, "could not create directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}
}
