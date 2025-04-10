package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	LogDir string
}

var (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
)

func disableColor() {
	Reset = ""
	Bold = ""
	Yellow = ""
	Green = ""
}

var config = Config{
	LogDir: ".logbook",
}

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	flag.Usage = showHelp
	flag.Parse()

	args := flag.Args()

	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			switch arg {
			case "-h", "--help":
				showHelp()
				return
			case "--version":
				fmt.Println("logbook version 1.0.0")
				return
			}
		}
	}

	logDir := filepath.Join(currentDir, config.LogDir)
	logFile := filepath.Join(logDir, time.Now().Format("2006-01-02")+".md")

	cmd := ""
	if len(args) > 0 {
		cmd = args[0]
		args = args[1:]
	}

	switch cmd {
	case "", "edit":
		editEntry(logFile)
	case "add":
		appendToEntry(logFile, args)
	case "read":
		readEntry(logFile)
	case "ls":
		listEntries(logDir)
	case "grep":
		grepEntries(logDir, args)
	case "logfile":
		fmt.Println(logFile)
	case "logdir":
		fmt.Println(logDir)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		showHelp()
		os.Exit(1)
	}
}

// determine which editor to use
func defaultEditor() string {
	if os.Getenv("NO_COLOR") != "" {
		disableColor()
	}
	for _, arg := range os.Args {
		if arg == "--no-color" {
			disableColor()
		}
	}

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

	fmt.Fprintf(os.Stderr, "No suitable editor found. Set the EDITOR environment variable.\n")
	os.Exit(1)
	return ""
}

func editEntry(file string) {
	if _, err := os.Stat(filepath.Dir(file)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(file), 0755)
	}
	editor := defaultEditor()
	cmd := exec.Command(editor, file)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func appendToEntry(file string, lines []string) {
	if len(lines) == 0 {
		fmt.Fprintln(os.Stderr, "Nothing to add")
		return
	}
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %v\n", err)
		return
	}
	defer f.Close()
	_, _ = f.WriteString(strings.Join(lines, " ") + "\n")
	fmt.Printf("added entry to %s\n", filepath.Base(file))
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
	cmd.Run()
}

func listEntries(dir string) {
	files, _ := filepath.Glob(filepath.Join(dir, "*.md"))
	for _, f := range files {
		fmt.Println(filepath.Base(f))
	}
}

func grepEntries(dir string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "nothing to grep")
		return
	}
	cmd := exec.Command("grep", append([]string{"-iR", "--color"}, append(args, dir)...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func showHelp() {
	fmt.Fprintf(os.Stdout, `%[2]susage:%[1]s logbook [options]

%[2]sexamples:%[1]s
  %[3]slogbook%[1]s edit
  %[3]slogbook%[1]s add "new logbook entry"
  %[3]slogbook%[1]s read
  %[3]slogbook%[1]s ls
  %[3]slogbook%[1]s grep "search term"

%[2]soptions:%[1]s
  %[3]s-h, --help%[1]s           show this help message and exit
  %[3]s--version%[1]s            show version and exit

`, Reset, Bold, Yellow)
}