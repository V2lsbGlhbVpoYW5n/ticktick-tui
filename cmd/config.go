package cmd

import (
	"fmt"
	"os"

	"ticktick-tui/internal/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置管理命令",
	Long:  `管理TickTick CLI配置，包括OAuth凭据和其他设置。`,
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "设置配置值",
	Long:  `设置配置键值对。包括client_id、client_secret、redirect_uri。`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if err := core.SaveConfig(key, value); err != nil {
			fmt.Fprintf(os.Stderr, "无法保存配置：%v\n", err)
			os.Exit(1)
		}

		fmt.Printf("配置已设置：%s = %s\n", key, value)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有配置",
	Long:  `列出当前所有的配置键值对。`,
	Run: func(cmd *cobra.Command, args []string) {
		settings := viper.AllSettings()

		if len(settings) == 0 {
			fmt.Println("没有找到配置")
			return
		}

		fmt.Println("当前配置：")
		for key, value := range settings {
			// 隐藏敏感信息
			if key == "access_token" || key == "client_secret" {
				fmt.Printf("%s = ***\n", key)
			} else {
				fmt.Printf("%s = %v\n", key, value)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(listCmd)
}
