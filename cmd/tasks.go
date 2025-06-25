package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"ticktick-tui/internal/client"
	"ticktick-tui/internal/models"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "任务管理命令",
	Long:  `管理TickTick任务，包括创建、更新、删除和查看任务。`,
}

var getTaskCmd = &cobra.Command{
	Use:   "get <project_id> <task_id>",
	Short: "获取指定任务",
	Long:  `根据项目ID和任务ID获取任务详情。`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		projectID := args[0]
		taskID := args[1]

		task, err := client.GetTask(projectID, taskID)
		if err != nil {
			fmt.Printf("获取任务失败：%v\n", err)
			os.Exit(1)
		}

		printTaskJSON(task)
	},
}

var createTaskCmd = &cobra.Command{
	Use:   "create",
	Short: "创建新任务",
	Long:  `创建一个新的TickTick任务。`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()

		title, _ := cmd.Flags().GetString("title")
		projectID, _ := cmd.Flags().GetString("project")
		content, _ := cmd.Flags().GetString("content")
		desc, _ := cmd.Flags().GetString("desc")
		priority, _ := cmd.Flags().GetInt("priority")

		if title == "" {
			fmt.Println("错误：任务标题不能为空")
			os.Exit(1)
		}

		if projectID == "" {
			fmt.Println("错误：项目ID不能为空")
			os.Exit(1)
		}

		task := &models.Task{
			Title:     title,
			ProjectID: projectID,
			Content:   content,
			Desc:      desc,
			Priority:  models.TaskPriority(priority),
		}

		// 处理日期参数
		if dueDate, _ := cmd.Flags().GetString("due"); dueDate != "" {
			if parsed, err := time.Parse("2006-01-02", dueDate); err == nil {
				task.DueDate = &models.TickTickTime{Time: parsed}
			} else {
				fmt.Printf("无效的截止日期格式：%s（应为YYYY-MM-DD）\n", dueDate)
				os.Exit(1)
			}
		}

		createdTask, err := client.CreateTask(task)
		if err != nil {
			fmt.Printf("创建任务失败：%v\n", err)
			os.Exit(1)
		}

		fmt.Println("任务创建成功：")
		printTaskJSON(createdTask)
	},
}

var updateTaskCmd = &cobra.Command{
	Use:   "update <task_id>",
	Short: "更新任务",
	Long:  `更新指定的任务信息。`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		taskID := args[0]

		projectID, _ := cmd.Flags().GetString("project")
		if projectID == "" {
			fmt.Println("错误：项目ID不能为空")
			os.Exit(1)
		}

		task := &models.Task{
			ID:        taskID,
			ProjectID: projectID,
		}

		// 更新字段
		if title, _ := cmd.Flags().GetString("title"); title != "" {
			task.Title = title
		}
		if content, _ := cmd.Flags().GetString("content"); content != "" {
			task.Content = content
		}
		if desc, _ := cmd.Flags().GetString("desc"); desc != "" {
			task.Desc = desc
		}
		if priority, _ := cmd.Flags().GetInt("priority"); cmd.Flags().Changed("priority") {
			if priority != 0 && priority != 1 && priority != 3 && priority != 5 {
				fmt.Println("错误：任务优先级必须为0，1，3，5中的一个")
				os.Exit(1)
			}
			task.Priority = models.TaskPriority(priority)
		}

		updatedTask, err := client.UpdateTask(taskID, task)
		if err != nil {
			fmt.Printf("更新任务失败：%v\n", err)
			os.Exit(1)
		}

		fmt.Println("任务更新成功：")
		printTaskJSON(updatedTask)
	},
}

var completeTaskCmd = &cobra.Command{
	Use:   "complete <project_id> <task_id>",
	Short: "完成任务",
	Long:  `标记指定任务为已完成。`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		projectID := args[0]
		taskID := args[1]

		err := client.CompleteTask(projectID, taskID)
		if err != nil {
			fmt.Printf("完成任务失败：%v\n", err)
			os.Exit(1)
		}

		fmt.Println("任务已标记为完成")
	},
}

var deleteTaskCmd = &cobra.Command{
	Use:   "delete <project_id> <task_id>",
	Short: "删除任务",
	Long:  `删除指定的任务。`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		projectID := args[0]
		taskID := args[1]

		err := client.DeleteTask(projectID, taskID)
		if err != nil {
			fmt.Printf("删除任务失败：%v\n", err)
			os.Exit(1)
		}

		fmt.Println("任务删除成功")
	},
}

func getClient() *client.Client {
	token := viper.GetString("access_token")
	if token == "" {
		token = viper.GetString("token")
	}

	if token == "" {
		fmt.Println("错误：未找到访问令牌")
		fmt.Println("请先运行 'ticktick-tui auth login' 进行身份验证")
		os.Exit(1)
	}

	return client.NewClient(token)
}

func printTaskJSON(task *models.Task) {
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		fmt.Printf("格式化输出失败：%v\n", err)
		return
	}
	fmt.Println(string(data))
}

func init() {
	rootCmd.AddCommand(tasksCmd)
	tasksCmd.AddCommand(getTaskCmd)
	tasksCmd.AddCommand(createTaskCmd)
	tasksCmd.AddCommand(updateTaskCmd)
	tasksCmd.AddCommand(completeTaskCmd)
	tasksCmd.AddCommand(deleteTaskCmd)

	// 创建任务的标志
	createTaskCmd.Flags().StringP("title", "t", "", "任务标题（必需）")
	createTaskCmd.Flags().StringP("project", "p", "", "项目ID（必需）")
	createTaskCmd.Flags().StringP("content", "c", "", "任务内容")
	createTaskCmd.Flags().StringP("desc", "d", "", "任务描述")
	createTaskCmd.Flags().IntP("priority", "r", 0, "任务优先级（0 (Low)，1，3，5 (High)）")
	createTaskCmd.Flags().String("due", "", "截止日期（YYYY-MM-DD格式）")

	// 更新任务的标志
	updateTaskCmd.Flags().StringP("title", "t", "", "任务标题")
	updateTaskCmd.Flags().StringP("project", "p", "", "项目ID（必需）")
	updateTaskCmd.Flags().StringP("content", "c", "", "任务内容")
	updateTaskCmd.Flags().StringP("desc", "d", "", "任务描述")
	updateTaskCmd.Flags().IntP("priority", "r", 0, "任务优先级（0 (Low)，1，3，5 (High)）")
	updateTaskCmd.MarkFlagRequired("project")
}
