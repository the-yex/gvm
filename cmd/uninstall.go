/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/the-yex/gvm/internal/core"
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

		err := core.UninstallVersion(args[0])
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		cmd.Printf("Uninstalled %s successfully\n", args[0])
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
