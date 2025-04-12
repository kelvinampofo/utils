package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	customDir  string
	useVanilla bool
)

var rootCmd = &cobra.Command{
	Use:   "new-sketch [app-name]",
	Short: "Scaffold a new Vite app for rapid UI prototyping",
	Long: `new-sketch creates a new Vite app in your sketches directory for rapid UI prototyping.

Examples:
  new-sketch my-app                 # uses React template by default
  new-sketch --vanilla my-app       # uses vanilla TS template
  new-sketch                        # uses today's date as the app name
  new-sketch -d ~/projects my-app

For feedback or issues, visit: https://github.com/kelvinampofo/utils`,
	Run: runNewSketch,
}

// runNewSketch scaffolds a new Vite app in the sketches directory.
func runNewSketch(cmd *cobra.Command, args []string) {
	// default sketch directory
	sandboxDir := filepath.Join(os.Getenv("HOME"), "Developer", "workspaces", "sketches")

	if customDir != "" {
		sandboxDir = customDir
	}

	// fallback to today's date if no app name is given
	var appName string
	if len(args) > 0 {
		appName = args[0]
	} else {
		appName = time.Now().Format("2006-01-02")
	}

	dest := filepath.Join(sandboxDir, appName)

	// ensure sandbox directory exists
	if err := os.MkdirAll(sandboxDir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Error creating sketch directory:", err)
		os.Exit(1)
	}

	fmt.Println("Creating Vite app:", appName)

	// select template (react-ts or vanilla-ts)
	template := "react-ts"
	if useVanilla {
		template = "vanilla-ts"
	}

	// run vite scaffolding in sketches dir
	createCmd := exec.Command("npm", "create", "vite@latest", appName, "--", "--template", template)
	createCmd.Dir = sandboxDir
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr

	if err := createCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running create command: %v\n", err)
		os.Exit(1)
	}

	// post-init: install dependencies + open in VS Code
	postCmds := [][]string{
		{"npm", "i"},
		{"code", "."},
	}

	for _, commandArgs := range postCmds {
		command := exec.Command(commandArgs[0], commandArgs[1:]...)
		command.Dir = dest
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running %v: %v\n", commandArgs, err)
			os.Exit(1)
		}
	}
}

func init() {
	rootCmd.Flags().StringVarP(&customDir, "dir", "d", "", "Custom sandbox directory (default: ~/Developer/workspaces/sketches)")
	rootCmd.Flags().BoolVar(&useVanilla, "vanilla", false, "Scaffold using Vanilla TS instead of React")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
