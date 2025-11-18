/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/spf13/cobra"
	"github.com/the-yex/gvm/internal/consts"
	list2 "github.com/the-yex/gvm/internal/tui/list"
	"github.com/the-yex/gvm/pkg"
	"slices"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Go versions",
	Long: `Display all Go versions.
Example:
  gvm list
    Show all Go versions installed locally.
  gvm list -r
    Show all available Go versions remotely.`,
	Aliases: []string{"l", "ls"},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		remote, _ := cmd.Flags().GetBool("remote")
		vk := consts.All
		if remote {
			kind, _ := cmd.Flags().GetString("type")
			vk, err = consts.ParseVersionKind(kind)
			if err != nil {
				return err
			}
		}

		versions, err := pkg.NewVManager(remote, pkg.WithLocal()).List(vk)
		if err != nil {
			return err
		}
		versions = slices.Compact(versions)
		items := make([]list.Item, len(versions))
		for index, v := range versions {
			items[index] = v
		}
		title := list2.LOCAL
		if remote {
			title = list2.Remote
		}
		list2.NewListProgram(items, title).Run()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("remote", "r", false, "List remote Go versions")
	listCmd.Flags().StringP("type", "t", string(consts.All), "Version type (default all): stable | unstable | archived ")
}
