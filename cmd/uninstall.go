/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/the-yex/gvm/internal/core"
	"github.com/the-yex/gvm/pkg"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [version]",
	Short: "Uninstall a specific Go version",
	Long: `Remove an installed Go version from your local environment.

Examples:
  gvm uninstall 1.21.0
  gvm uninstall 1.20.5

This command will delete the corresponding Go version directory
from your installation path.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		v := pkg.LocalInstalled(version)
		if v == nil {
			cmd.Printf("Version %q is not installed. Install it with \"gvm install %s\" first.\n", version, version)
			return
		}
		err := core.UninstallVersion(v.LocalDir())
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		cmd.Printf("Uninstalled %s successfully\n", version)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uninstallCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uninstallCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
