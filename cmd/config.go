/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/the-yex/gvm/internal/consts"
	"github.com/the-yex/gvm/internal/prettyout"
	"github.com/the-yex/gvm/internal/utils"
	"os"
	"reflect"
	"slices"
	"strings"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage gvm configuration",
		Long: `View and modify gvm configuration.

Examples:
  gvm config list           # show all config
  gvm config set mirrors https://mirrors.aliyun.com/golang/
  gvm config get mirrors    # show a specific config
  gvm config unset mirrors  # remove a config item`,
	}
	configListCmd = &cobra.Command{
		Use:     "list",
		Short:   "List all configuration values",
		Aliases: []string{"l", "ls"},
		Run: func(cmd *cobra.Command, args []string) {
			settings := viper.AllSettings()
			for k, v := range settings {
				fmt.Printf("%s: %v\n", k, v)
			}
		},
	}
	configGetCmd = &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration value",
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			val := viper.Get(key)
			if val == nil {
				return fmt.Errorf(
					"invalid key: [%s]\nSupported keys:\n- %s",
					key,
					strings.Join(viper.AllKeys(), "\n- "),
				)
			}
			fmt.Printf("%s: %v\n", key, val)
			return nil
		},
		Args: cobra.MinimumNArgs(1),
	}
	configSetCmd = &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration value",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			setConfig(args...)
		},
	}
	configUnsetCmd = &cobra.Command{
		Use:   "unset [key]",
		Short: "Remove a configuration value",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			unSetConfig(args...)
		},
	}
)

/*
*
规则：

 1. 如果是普通值，直接删除整个 key
    gvm config unset goroots   ❌（报错，因为 goroots 是 list）
    gvm config unset mirror    ✅（删除整个 key）

 2. 如果是 list，必须带额外参数（要删除的值）
    gvm config unset goroots /Users/zouyuxi/sdk
*/
func unSetConfig(args ...string) {
	key := args[0]
	existValues := viper.Get(key)
	if existValues == nil {
		fmt.Printf("%s not set\n", key)
		return
	}
	kind := reflect.TypeOf(existValues).Kind()
	switch kind {
	case reflect.Slice:
		// list 必须至少有 2 个参数
		if len(args) < 2 {
			prettyout.PrettyWarm(os.Stdout, "%s is a list, you must specify a value to remove\n", key)
			return
		}
		toRemove := args[1:]
		oldVals := viper.GetStringSlice(key)
		newVals := []string{}
		for _, v := range oldVals {
			if !slices.Contains(toRemove, v) {
				newVals = append(newVals, v)
			}
		}
		if len(newVals) == 0 {
			newVals = append(newVals, consts.EMPTY_INFO)
		}
		viper.Set(key, newVals)
	default:
		viper.Set(key, "")
	}
	viper.WriteConfig()
}

func setConfig(args ...string) {
	key := args[0]
	values := args[1:]
	existValues := viper.Get(key)
	if existValues == nil {
		prettyout.PrettyWarm(os.Stdout, "%s\n", "no such config key support")
		return
	}
	if reflect.TypeOf(existValues).Kind() == reflect.Slice {
		newValues := append(viper.GetStringSlice(key), values...)
		newValues = utils.Unique(newValues)
		newValues = slices.DeleteFunc(newValues, func(s string) bool {
			return s == consts.EMPTY_INFO
		})
		viper.Set(key, newValues)
	} else {
		viper.Set(key, values[0])
	}

	viper.WriteConfig()
}
func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
