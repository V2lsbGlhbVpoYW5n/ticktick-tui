package core

import (
	"fmt"
	"ticktick-tui/internal/auth"
	"ticktick-tui/internal/client"
	"ticktick-tui/internal/models"

	"github.com/spf13/viper"
)

func SaveConfig(key, value string) error {

	viper.Set(key, nil)
	viper.Set(key, value)

	if err := viper.WriteConfig(); err != nil {
		if err := viper.SafeWriteConfig(); err != nil {
			return err
		}
	}

	return nil
}

func GetAuthURL() string {
	client := &auth.OAuthClient{
		ClientID:    viper.GetString("client_id"),
		RedirectURI: viper.GetString("redirect_uri"),
	}

	return client.GetAuthURL()
}

func GetToken(code string) (*models.OAuthToken, error) {
	client := &auth.OAuthClient{
		ClientID:     viper.GetString("client_id"),
		ClientSecret: viper.GetString("client_secret"),
		RedirectURI:  viper.GetString("redirect_uri"),
	}
	scope := "tasks:read tasks:write"

	return client.ExchangeCodeForToken(code, scope)
}

func getClient() *client.Client {
	token := viper.GetString("access_token")
	if token == "" {
		token = viper.GetString("token")
	}

	if token == "" {
		return nil
	}

	return client.NewClient(token)
}

func GetProjects() ([]models.Project, error) {
	client := getClient()
	if client == nil {
		return nil, fmt.Errorf("未找到访问令牌")
	}

	projects, err := client.GetProjects()
	if err != nil {
		return nil, fmt.Errorf("获取项目列表失败：%v", err)
	}

	return projects, nil
}

func GetTasks(projectID string) ([]models.Task, error) {
	client := getClient()
	if client == nil {
		return nil, fmt.Errorf("未找到访问令牌")
	}

	tasks, err := client.GetProjectData(projectID)
	if err != nil {
		return nil, fmt.Errorf("获取任务列表失败：%v", err)
	}

	return tasks.Tasks, nil
}
