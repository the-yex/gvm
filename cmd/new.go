/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/the-yex/gvm/pkg"
	"os"
	"os/exec"
	"path/filepath"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go project with the current active version",
	Long: `Create a new Go project initialized with the Go version currently set by gvm.

Example:
  gvm new myapp
This will create a folder 'myapp', initialize a Go module,
and set it up using the active Go version.`,
	Args: cobra.ExactArgs(1), // require project name
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		version, _ := cmd.Flags().GetString("version")
		module, _ := cmd.Flags().GetString("module")
		v := pkg.LocalInstalled(version)
		if version != "" && v == nil {
			cmd.PrintErrf("❌ Go version %s has not been installed\n", version)
			return
		}
		if module == "" {
			module = projectName
		}
		if _, err := os.Stat(projectName); err == nil {
			cmd.PrintErrf("❌ Project folder '%s' already exists\n", projectName)
			return
		}
		if err := os.Mkdir(projectName, 0755); err != nil {
			cmd.PrintErrf("❌ Failed to create project folder: %v\n", err)
			return
		}

		cmdName := "go"
		if v != nil {
			cmdName = filepath.Join(v.LocalDir(), "bin", "go")
		}
		cmdStr := exec.Command(cmdName, "mod", "init", module)

		cmdStr.Dir = projectName
		if output, err := cmdStr.CombinedOutput(); err != nil {
			cmd.PrintErrf("❌ Failed to initialize go.mod: %s\n%s\n", err, string(output))
			return
		}
		mainGo := fmt.Sprintf(`package main

import "fmt"

func main() {
    fmt.Println("Hello from %s!")
}
`, projectName)
		if err := os.WriteFile(filepath.Join(projectName, "main.go"), []byte(mainGo), 0644); err != nil {
			cmd.PrintErrf("❌ Failed to write main.go: %v\n", err)
			return
		}

		cmd.Printf("✅ New Go project '%s' created successfully!\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("version", "V", "", "Go version(default in current version)")
	newCmd.Flags().StringP("module", "m", "", "Go module path (default is project name)")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
