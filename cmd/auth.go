package cmd

import (
	"fmt"
	"os"
	"ticktick-tui/internal/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "OAuth认证相关命令",
	Long:  `管理TickTick OAuth认证流程，包括生成授权URL和交换访问令牌。`,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "生成OAuth授权URL",
	Long:  `生成OAuth授权URL，用户需要在浏览器中打开此URL进行授权。`,
	Run: func(cmd *cobra.Command, args []string) {
		clientID := viper.GetString("client_id")
		redirectURI := viper.GetString("redirect_uri")

		if clientID == "" || redirectURI == "" {
			fmt.Println("错误：请先配置client_id和redirect_uri")
			os.Exit(1)
		}

		authURL := core.GetAuthURL()
		fmt.Println("请在浏览器中打开以下URL进行授权：")
		fmt.Println(authURL)
		fmt.Println("\n授权完成后，从重定向URL中复制授权码，然后运行：")
		fmt.Println("ticktick-tui auth token <authorization_code>")
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token <authorization_code>",
	Short: "使用授权码获取访问令牌",
	Long:  `使用从OAuth流程获得的授权码来获取访问令牌。`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clientID := viper.GetString("client_id")
		clientSecret := viper.GetString("client_secret")
		redirectURI := viper.GetString("redirect_uri")

		if clientID == "" || clientSecret == "" || redirectURI == "" {
			fmt.Println("错误：请先配置client_id、client_secret和redirect_uri")
			os.Exit(1)
		}

		code := args[0]

		token, err := core.GetToken(code)
		if err != nil {
			fmt.Printf("获取访问令牌失败：%v\n", err)
			os.Exit(1)
		}

		viper.Set("access_token", token.AccessToken)
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("保存配置失败：%v\n", err)
			os.Exit(1)
		}

		fmt.Println("访问令牌获取成功并已保存到配置文件！")
		fmt.Printf("访问令牌：%s\n", token.AccessToken)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(tokenCmd)
}
