/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/viper"
	"github.com/the-yex/gvm/internal/consts"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func Root() *cobra.Command { return rootCmd }

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gvm",
	Short: "Go version manager for installing and switching between multiple Go versions",
	Long: `gvm is a Go version management tool, similar to nvm for Node.js.

With gvm you can:
- Install specific versions of Go
- Switch between multiple installed Go versions
- Uninstall versions you no longer need

Perfect for developers who work on projects that require different Go versions.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initializeConfig()
	},
	Version: consts.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var verbose bool

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gvm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initializeConfig() {
	viper.SetEnvPrefix("gvm")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	viper.AutomaticEnv()
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home + "/.gvm")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		// basic configs
		viper.Set(consts.CONFIG_MIRROR, consts.DEFAULT_MIRROR)
		viper.Set(consts.CONFIG_GOROOT, []string{home + consts.DEFAULT_GOROOT})
		viper.SafeWriteConfig()
	}
}
