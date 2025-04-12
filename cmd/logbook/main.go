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
	FilePermission:  0755, // → rwxr-xr-x,
	EntryPermission: 0644, // → rw-r--r--
}

// editCmd opens today's log or a specific date file
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

// addCmd appends text to today's log
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

// rootCmd is the base command called without any subcommands, e.g. "logbook"
var rootCmd = &cobra.Command{
	Use:   "logbook [subcommand]",
	Short: "A command-line interface for daily Markdown logbooks.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// register subcommands on the root
	rootCmd.AddCommand(
		addCmd,
		editCmd,
		grepCmd,
		logdirCmd,
		logfileCmd,
		lsCmd,
		readCmd,
	)
}

func main() {
	// disable completion option
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// if "logbook" is run with no subcommand and no flags,
	// Cobra will call rootCmd.Run and then exit.
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// resolveLogDir returns the full path to .logbook in the current directory
func resolveLogDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting current dir: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(currentDir, config.LogDir)
}

// resolveLogFile returns the full path to today's .md file
func resolveLogFile(dateStr string) string {
	// if dateStr is empty, use today's date
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	} else {
		_, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid date format: %s (expected YYYY-MM-DD)\n", dateStr)
			os.Exit(1)
		}
	}
	return filepath.Join(resolveLogDir(), dateStr+config.FileExtension)
}

func editEntry(file string) {
	path := filepath.Dir(file)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// create .logbook if it doesn't exist
		if mkErr := os.MkdirAll(path, config.FilePermission); mkErr != nil {
			fmt.Fprintf(os.Stderr, "could not create log directory: %v\n", mkErr)
			os.Exit(1)
		}
	}
	editor := defaultEditor()
	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runErr := cmd.Run(); runErr != nil {
		fmt.Fprintf(os.Stderr, "failed to open editor: %v\n", runErr)
		os.Exit(1)
	}
}

func appendToEntry(file string, lines []string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, config.EntryPermission)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	_, _ = f.WriteString(strings.Join(lines, " ") + "\n")
	fmt.Fprintf(os.Stdout, "added entry to \"%s\"\n", filepath.Base(file))
}

func readEntry(file string) {
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "less"
	}
	cmd := exec.Command(pager, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runErr := cmd.Run(); runErr != nil {
		fmt.Fprintf(os.Stderr, "failed to open pager: %v\n", runErr)
		os.Exit(1)
	}
}

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

func grepEntries(dir string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "nothing to grep.")
		os.Exit(1)
	}
	// prepend grep options, then append dir at the end
	cmd := exec.Command("grep", append([]string{"-iR", "--color"}, append(args, dir)...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runErr := cmd.Run(); runErr != nil {
		fmt.Fprintf(os.Stderr, "grep failed: %v\n", runErr)
		os.Exit(1)
	}
}

func defaultEditor() string {
	// respect EDITOR environment variable
	if editorEnv := os.Getenv("EDITOR"); editorEnv != "" {
		return editorEnv
	}

	editors := []string{"nvim", "vim", "nano", "emacs"}
	for _, editor := range editors {
		if _, err := exec.LookPath(editor); err == nil {
			return editor
		}
	}

	fmt.Fprintf(os.Stderr, "no suitable editor found. Set the EDITOR environment variable.\n")
	os.Exit(1)
	return ""
}
