/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/the-yex/gvm/internal/consts"
	"github.com/the-yex/gvm/pkg"
	"os"
	"regexp"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Switch to a specific Go version",
	Args:  cobra.MaximumNArgs(1),
	Long: `Switch the current Go environment to the specified version.

Examples:
  gvm use go1.21    # activate Go 1.21
  gvm use           # use version from go.mod (if available)`,
	Run: func(cmd *cobra.Command, args []string) {
		version := ""
		if len(args) == 0 {
			modVersion, err := getModVersion()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			version = modVersion
		} else {
			version = args[0]
		}
		localVersions, _ := pkg.NewVManager(false).List(consts.All)
		for _, localVersion := range localVersions {
			if localVersion.String() == version {
				if err := pkg.SwitchVersion(localVersion.LocalDir()); err != nil {
					cmd.Println(err.Error())
				}
				return
			}
		}
		cmd.Printf("Version %q not found. use  \"gvm install %s\" first\n", version, version)
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var goModReg = regexp.MustCompile(`(?m)^go\s+(\d+\.\d+(?:\.\d+)?(?:beta\d+|rc\d+)?)\s*(?:$|//.*)`)

func getModVersion() (string, error) {
	// Uses go.mod if available and version is omitted
	goModData, err := os.ReadFile("go.mod")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", errors.New("no go.mod file found")
		}
		return "", err
	}
	match := goModReg.FindStringSubmatch(string(goModData))
	if len(match) > 1 {
		return match[1], nil
	}
	return "", errors.New("no version found in go.mod")
}
