package tui

import (
	"fmt"
	"strings"
	"ticktick-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) renderStatusBar() string {
	// TODO: ADD project title to status bar in TaskListView
	if m.width == 0 {
		return ""
	}
	// Left section: Current view/mode, with fixed widths
	var leftSection string
	var mode string

	const leftSectionWidth = 10
	const modeWidth = 8

	switch m.state.CurrentView {
	case models.ConfigView:
		leftSection = statusLeftStyle.
			Foreground(BLACK).
			Background(BLUE).
			Width(leftSectionWidth).
			Render("CONFIG")
		mode = statusModeStyle.Width(modeWidth).Render("NORMAL")
	case models.AuthView:
		leftSection = statusLeftStyle.
			Foreground(BLACK).
			Background(RED).
			Width(leftSectionWidth).
			Render("AUTH")
		mode = statusModeStyle.Width(modeWidth).Render("INPUT")
	case models.ProjectListView:
		leftSection = statusLeftStyle.
			Foreground(BLACK).
			Background(BLUE).
			Width(leftSectionWidth).
			Render("PROJECTS")
		mode = statusModeStyle.Width(modeWidth).Render("NORMAL")
	case models.TaskListView:
		leftSection = statusLeftStyle.
			Foreground(BLACK).
			Background(BLUE).
			Width(leftSectionWidth).
			Render("TASKS")
		mode = statusModeStyle.Width(modeWidth).Render("NORMAL")
	case models.TaskDetailView:
		leftSection = statusLeftStyle.
			Foreground(BLACK).
			Background(BLUE).
			Width(leftSectionWidth).
			Render("TASK")
		mode = statusModeStyle.Width(modeWidth).Render("NORMAL")
	default:
		leftSection = statusLeftStyle.
			Foreground(BLACK).
			Background(BLUE).
			Width(leftSectionWidth).
			Render(" ")
		mode = statusModeStyle.Width(modeWidth).Render("NORMAL")
	}

	// Middle section: Messages or errors
	var middleSection string
	if m.state.Error != "" {
		middleSection = statusErrorStyle.Render("ERROR")
	} else if m.state.Loading {
		middleSection = statusMessageStyle.Render(fmt.Sprintf(" %s Loading... ", m.spinner.View()))
	}

	// Right section: Position info
	var rightSection string
	switch m.state.CurrentView {
	case models.ProjectListView:
		if len(m.state.Projects) > 0 {
			rightSection = statusRightStyle.Render(fmt.Sprintf("%d/%d", m.state.SelectedIndex+1, len(m.state.Projects)))
		}
	case models.TaskListView:
		if len(m.state.Tasks) > 0 {
			rightSection = statusRightStyle.Render(fmt.Sprintf("%d/%d", m.state.SelectedIndex+1, len(m.state.Tasks)))
		}
	case models.ConfigView:
		rightSection = statusRightStyle.Render(fmt.Sprintf("Field %d/3", m.state.SelectedIndex+1))
	}

	// Calculate available space for middle section
	leftWidth := lipgloss.Width(leftSection + mode)
	rightWidth := lipgloss.Width(rightSection)
	middleWidth := m.width - leftWidth - rightWidth

	// Truncate middle section if too long
	if middleWidth > 0 && lipgloss.Width(middleSection) > middleWidth {
		maxLen := middleWidth - 3 // Reserve space for "..."
		if maxLen > 0 {
			truncated := []rune(middleSection)
			if len(truncated) > maxLen {
				middleSection = string(truncated[:maxLen]) + "..."
			}
		}
	}

	// Fill remaining space
	remaining := m.width - lipgloss.Width(leftSection+mode+middleSection+rightSection)
	if remaining < 0 {
		remaining = 0
	}
	filler := strings.Repeat(" ", remaining)

	// Combine all sections
	statusContent := leftSection + mode + middleSection + filler + rightSection

	return statusBarStyle.Width(m.width).Render(statusContent)
}

func (m *Model) renderConfigForm() string {
	title := formTitleStyle.Render("TickTick 配置")

	var form strings.Builder

	for i, item := range m.state.CurrentItems {
		switch v := item.(type) {
		case textinput.Model:
			if i == m.state.SelectedIndex {
				v.TextStyle = formFocusedStyle
				v.PromptStyle = formFocusedStyle
				v.PlaceholderStyle = formFocusedPlaceHolderStyle
				v.Cursor.Style = formFocusedStyle
				form.WriteString(v.View())
			} else {
				v.TextStyle = formBlurredStyle
				v.PromptStyle = formBlurredStyle
				v.PlaceholderStyle = formBlurredPlaceHolderStyle
				form.WriteString(v.View())
			}
			form.WriteString("\n\n")
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		form.String(),
	)

	return formStyle.Width(m.width - 8).Render(content)
}

func (m *Model) renderAuthForm() string {
	title := formTitleStyle.Render("TickTick 授权")

	var form strings.Builder

	// Auth URL display
	if m.state.AuthURL != "" {
		form.WriteString(formBlurredStyle.Render("请在浏览器中打开以下链接进行授权:"))
		form.WriteString("\n")
		form.WriteString(formBlurredStyle.Render("已复制到剪贴板"))
		form.WriteString("\n")
		form.WriteString(authURLStyle.Width(m.width - 16).Render(m.state.AuthURL))
		form.WriteString("\n\n")
	}

	// Auth code input
	for i, item := range m.state.CurrentItems {
		switch v := item.(type) {
		case textinput.Model:
			if i == m.state.SelectedIndex {
				v.TextStyle = formFocusedStyle
				v.PromptStyle = formFocusedStyle
				v.PlaceholderStyle = formFocusedPlaceHolderStyle
				v.Cursor.Style = formFocusedStyle
				form.WriteString(v.View())
			} else {
				v.TextStyle = formBlurredStyle
				v.PromptStyle = formBlurredStyle
				v.PlaceholderStyle = formBlurredPlaceHolderStyle
				form.WriteString(v.View())
			}
			form.WriteString("\n\n")
		}
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		form.String(),
	)

	return formStyle.Width(m.width - 8).Render(content)
}

func (m *Model) renderProjectList() string {
	if len(m.state.Projects) == 0 {
		return lipgloss.NewStyle().
			Width(m.width).
			Padding(2, 2).
			Render("No projects found")
	}

	items := make([]list.Item, len(m.state.Projects))
	for i, project := range m.state.Projects {
		desc := "Active project"
		if project.Kind != "" {
			desc = project.Kind
		}
		if project.Closed {
			desc = "Archived"
		}
		if project.GroupID != "" {
			groupID := project.GroupID
			if len(groupID) > 4 {
				groupID = groupID[len(groupID)-4:]
			}
			desc += " • Group: " + groupID
		}

		title := project.Name
		if project.Color != "" {
			colorIndicator := lipgloss.NewStyle().Foreground(lipgloss.Color(project.Color)).Render("●")
			title = project.Name + " " + colorIndicator
		}

		items[i] = projectItem{
			title: title,
			desc:  desc,
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = listNormalTitleStyle
	delegate.Styles.NormalDesc = listNormalDescStyle
	delegate.Styles.SelectedTitle = listSelectedTitleStyle
	delegate.Styles.SelectedDesc = listSelectedDescStyle

	// Reset the list height to fit the current model height
	l := list.New(items, delegate, m.width-8, m.height-8)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	l.Paginator.Type = paginator.Dots
	l.Paginator.ActiveDot = paginatorActive
	l.Paginator.InactiveDot = paginatorInactive

	// Set selected index
	if m.state.SelectedIndex < len(items) {
		l.Select(m.state.SelectedIndex)
	}

	return listStyle.
		Width(m.width - 8).
		Render(l.View())
}

// Add this type definition for list items
type projectItem struct {
	title, desc string
}

func (i projectItem) Title() string       { return i.title }
func (i projectItem) Description() string { return i.desc }
func (i projectItem) FilterValue() string { return i.title }

func (m *Model) renderTaskList() string {
	if len(m.state.Tasks) == 0 {
		return lipgloss.NewStyle().
			Width(m.width).
			Padding(2, 2).
			Render("No tasks found")
	}

	items := make([]list.Item, len(m.state.Tasks))
	for i, task := range m.state.Tasks {
		title := task.Title

		if task.Priority != models.PriorityNone {
			var priorityIndicator string
			switch task.Priority {
			case models.PriorityLow:
				priorityIndicator = priorityLow
			case models.PriorityMedium:
				priorityIndicator = priorityMedium
			case models.PriorityHigh:
				priorityIndicator = priorityHigh
			default:
				priorityIndicator = priorityNone
			}
			title = title + " " + priorityIndicator
		}

		var desc string

		if task.DueDate != nil {
			due := task.DueDate.String()
			if due != "" {
				desc = "Due: " + due
			}
		}

		if task.Content != "" {
			if desc != "" {
				desc += " • "
			}
			content := task.Content
			if len(content) > 10 {
				content = content[:10] + "..."
			}
			desc += content
		}

		items[i] = taskItem{
			title: title,
			desc:  desc,
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = listNormalTitleStyle
	delegate.Styles.NormalDesc = listNormalDescStyle
	delegate.Styles.SelectedTitle = listSelectedTitleStyle
	delegate.Styles.SelectedDesc = listSelectedDescStyle

	l := list.New(items, delegate, m.width-8, m.height-8)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	l.Paginator.Type = paginator.Dots
	l.Paginator.ActiveDot = paginatorActive
	l.Paginator.InactiveDot = paginatorInactive

	// Set selected index
	if m.state.SelectedIndex < len(items) {
		l.Select(m.state.SelectedIndex)
	}

	return listStyle.
		Width(m.width - 8).
		Render(l.View())
}

type taskItem struct {
	title, desc string
}

func (i taskItem) Title() string       { return i.title }
func (i taskItem) Description() string { return i.desc }
func (i taskItem) FilterValue() string { return i.title }

func (m *Model) renderMessage() string {
	if m.state.Error != "" {
		return messageErrorStyle.Render(m.state.Error)
	}
	if m.state.Message != "" {
		return messageStyle.Render(m.state.Message)
	}
	return ""
}

func (m *Model) renderHelp() string {
	if m.width == 0 {
		return ""
	}

	var helpItems []string

	switch m.state.CurrentView {
	case models.ConfigView:
		helpItems = []string{
			m.helpKey("Ctrl+c", ""),
			m.helpKey("Up/Down", "Select"),
			m.helpKey("Enter", ""),
		}
	case models.AuthView:
		helpItems = []string{
			m.helpKey("Ctrl+c", "Exit"),
			m.helpKey("Enter", "Submit"),
		}
	case models.ProjectListView:
		helpItems = []string{
			m.helpKey("Ctrl+c", "Exit"),
			m.helpKey("Up/Down", "Select"),
			m.helpKey("Enter", "Open"),
			m.helpKey("a", "New"),
			m.helpKey("d", "Delete"),
			m.helpKey("e", "Edit"),
		}
	case models.TaskListView:
		helpItems = []string{
			m.helpKey("Ctrl+c", "Exit"),
			m.helpKey("Up/Down", "Select"),
			m.helpKey("Enter", "Open"),
			m.helpKey("a", "New"),
			m.helpKey("d", "Delete"),
			m.helpKey("e", "Edit"),
			m.helpKey("Space", "[Un]Complete"),
			m.helpKey("Esc", "Back"),
		}
	case models.TaskDetailView:
		helpItems = []string{
			m.helpKey("Ctrl+c", "Exit"),
			m.helpKey("e", "Edit"),
			m.helpKey("d", "Delete"),
			m.helpKey("Space", "[Un]Complete"),
			m.helpKey("Esc", "Back"),
		}
	case models.DeleteConfirmView:
		helpItems = []string{
			m.helpKey("Ctrl+c", "Exit"),
			m.helpKey("Enter", "Confirm"),
			m.helpKey("Esc", "Back"),
		}

	default:
		helpItems = []string{
			m.helpKey("Ctrl+c", "Exit"),
		}
	}

	// Join help items with separators
	helpContent := strings.Join(helpItems, "  •  ")

	// Truncate if too long
	if lipgloss.Width(helpContent) > m.width-4 {
		maxLen := m.width - 7 // Reserve space for "..."
		if maxLen > 0 {
			truncated := []rune(helpContent)
			if len(truncated) > maxLen {
				helpContent = string(truncated[:maxLen]) + "..."
			}
		}
	}

	return helpStyle.Width(m.width).Render(helpContent)
}

func (m *Model) helpKey(key, desc string) string {
	return helpKeyStyle.Render(key) + " " + helpDescStyle.Render(desc)
}
