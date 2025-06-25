package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"ticktick-tui/internal/models"

	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "项目管理命令",
	Long:  `管理TickTick项目，包括创建、更新、删除和查看项目。`,
}

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有项目",
	Long:  `获取用户的所有项目列表。`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()

		projects, err := client.GetProjects()
		if err != nil {
			fmt.Printf("获取项目列表失败：%v\n", err)
			os.Exit(1)
		}

		// Ensure projects is of type []models.Project
		var typedProjects []models.Project
		typedProjects = append(typedProjects, projects...)
		printProjectsJSON(typedProjects)
	},
}

var getProjectDataCmd = &cobra.Command{
	Use:   "data <project_id>",
	Short: "获取项目完整数据",
	Long:  `获取项目的完整数据，包括任务和列。`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		projectID := args[0]

		projectData, err := client.GetProjectData(projectID)
		if err != nil {
			fmt.Printf("获取项目数据失败：%v\n", err)
			os.Exit(1)
		}

		printProjectDataJSON(projectData)
	},
}

var createProjectCmd = &cobra.Command{
	Use:   "create",
	Short: "创建新项目",
	Long:  `创建一个新的TickTick项目。`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()

		name, _ := cmd.Flags().GetString("name")
		color, _ := cmd.Flags().GetString("color")
		viewMode, _ := cmd.Flags().GetString("view-mode")
		kind, _ := cmd.Flags().GetString("kind")

		if name == "" {
			fmt.Println("错误：项目名称不能为空")
			os.Exit(1)
		}

		project := &models.Project{
			Name:     name,
			Color:    color,
			ViewMode: viewMode,
			Kind:     kind,
		}

		createdProject, err := client.CreateProject(project)
		if err != nil {
			fmt.Printf("创建项目失败：%v\n", err)
			os.Exit(1)
		}

		fmt.Println("项目创建成功：")
		printProjectJSON(createdProject)
	},
}

var updateProjectCmd = &cobra.Command{
	Use:   "update <project_id>",
	Short: "更新项目",
	Long:  `更新指定的项目信息。`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		projectID := args[0]

		project := &models.Project{}

		if name, _ := cmd.Flags().GetString("name"); name != "" {
			project.Name = name
		}
		if color, _ := cmd.Flags().GetString("color"); color != "" {
			project.Color = color
		}
		if viewMode, _ := cmd.Flags().GetString("view-mode"); viewMode != "" {
			project.ViewMode = viewMode
		}
		if kind, _ := cmd.Flags().GetString("kind"); kind != "" {
			project.Kind = kind
		}

		updatedProject, err := client.UpdateProject(projectID, project)
		if err != nil {
			fmt.Printf("更新项目失败：%v\n", err)
			os.Exit(1)
		}

		fmt.Println("项目更新成功：")
		printProjectJSON(updatedProject)
	},
}

var deleteProjectCmd = &cobra.Command{
	Use:   "delete <project_id>",
	Short: "删除项目",
	Long:  `删除指定的项目。`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		projectID := args[0]

		err := client.DeleteProject(projectID)
		if err != nil {
			fmt.Printf("删除项目失败：%v\n", err)
			os.Exit(1)
		}

		fmt.Println("项目删除成功")
	},
}

func printProjectJSON(project *models.Project) {
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		fmt.Printf("格式化输出失败：%v\n", err)
		return
	}
	fmt.Println(string(data))
}

func printProjectsJSON(projects []models.Project) {
	// Sort projects by the Ordered field
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].SortOrder < projects[j].SortOrder
	})

	grouped := make(map[string][]models.Project)

	// Pickout Archived projects
	for _, p := range projects {
		if p.Closed {
			grouped["Archived"] = append(grouped["Archived"], p)
			continue
		}

		if p.GroupID != "" {
			grouped[p.GroupID] = append(grouped[p.GroupID], p)
		} else {
			grouped["Ungrouped"] = append(grouped["Ungrouped"], p)
		}
	}

	// Prepare ordered output: groups with normalId first, then Ungrouped, then Archived
	var ordered []interface{}
	// Add groups with normalId (i.e., not "Ungrouped" or "Archived")
	var groupIDs []string
	for groupID := range grouped {
		if groupID != "Ungrouped" && groupID != "Archived" {
			groupIDs = append(groupIDs, groupID)
		}
	}
	sort.Strings(groupIDs)
	for _, groupID := range groupIDs {
		ordered = append(ordered, map[string]interface{}{
			"groupId":  groupID,
			"projects": grouped[groupID],
		})
	}
	// Add Ungrouped
	if ungrouped, ok := grouped["Ungrouped"]; ok {
		ordered = append(ordered, map[string]interface{}{
			"groupId":  "Ungrouped",
			"projects": ungrouped,
		})
	}
	// Add Archived
	if archived, ok := grouped["Archived"]; ok {
		ordered = append(ordered, map[string]interface{}{
			"groupId":  "Archived",
			"projects": archived,
		})
	}

	data, err := json.MarshalIndent(ordered, "", "  ")

	if err != nil {
		fmt.Printf("格式化输出失败：%v\n", err)
		return
	}
	fmt.Println(string(data))
}

func printProjectDataJSON(projectData *models.ProjectData) {
	data, err := json.MarshalIndent(projectData, "", "  ")
	if err != nil {
		fmt.Printf("格式化输出失败：%v\n", err)
		return
	}
	fmt.Println(string(data))
}

func init() {
	rootCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(listProjectsCmd)
	projectsCmd.AddCommand(getProjectDataCmd)
	projectsCmd.AddCommand(createProjectCmd)
	projectsCmd.AddCommand(updateProjectCmd)
	projectsCmd.AddCommand(deleteProjectCmd)

	createProjectCmd.Flags().StringP("name", "n", "", "项目名称（必需）")
	createProjectCmd.Flags().StringP("color", "c", "", "项目颜色（如：#F18181）")
	createProjectCmd.Flags().String("view-mode", "list", "视图模式（list, kanban, timeline）")
	createProjectCmd.Flags().String("kind", "TASK", "项目类型（TASK, NOTE）")
	createProjectCmd.MarkFlagRequired("name")

	updateProjectCmd.Flags().StringP("name", "n", "", "项目名称")
	updateProjectCmd.Flags().StringP("color", "c", "", "项目颜色（如：#F18181）")
	updateProjectCmd.Flags().String("view-mode", "", "视图模式（list, kanban, timeline）")
	updateProjectCmd.Flags().String("kind", "", "项目类型（TASK, NOTE）")
}
