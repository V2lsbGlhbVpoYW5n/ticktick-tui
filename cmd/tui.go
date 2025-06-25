package cmd

import (
	"fmt"
	"os"
	"ticktick-tui/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "启动TUI界面",
	Long:  `启动交互式终端用户界面来管理TickTick任务和项目。`,
	Run: func(cmd *cobra.Command, args []string) {
		model := tui.NewModel()

		p := tea.NewProgram(model, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Printf("启动TUI失败: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
