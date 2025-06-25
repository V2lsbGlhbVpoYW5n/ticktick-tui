package tui

import (
	"os"
	"ticktick-tui/internal/client"
	"ticktick-tui/internal/models"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

type Model struct {
	client *client.Client
	state  *models.AppState
	width  int
	height int

	// UI Components
	spinner spinner.Model

	// Config Inputs
	configInputs []textinput.Model

	// Auth Inputs
	authInputs []textinput.Model
}

type (
	projectsLoadedMsg []models.Project
	tasksLoadedMsg    []models.Task

	taskCreatedMsg    *models.Task
	projectCreatedMsg *models.Project

	configSavedMsg    struct{}
	tokenExchangedMsg struct{}
)

func NewModel() *Model {
	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	// Initialize text inputs
	clientIDInput := textinput.New()
	clientIDInput.Placeholder = "Client ID"

	clientSecretInput := textinput.New()
	clientSecretInput.Placeholder = "Client Secret"

	redirectURIInput := textinput.New()
	redirectURIInput.Placeholder = "Redirect URI"

	authCodeInput := textinput.New()
	authCodeInput.Placeholder = "Authorization Code"

	// Initialize model
	m := &Model{
		spinner: s,
		state: &models.AppState{
			Projects:      []models.Project{},
			Tasks:         []models.Task{},
			CurrentItems:  []any{},
			SelectedIndex: 0,
			Loading:       false,
		},
		configInputs: []textinput.Model{
			clientIDInput,
			clientSecretInput,
			redirectURIInput,
		},
		authInputs: []textinput.Model{
			authCodeInput,
		},
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	// Set window title
	cmds = append(cmds, tea.SetWindowTitle("ticktick-tui"))

	// Set initial view
	// Check config and set initial view
	if viper.GetString("client_id") == "" || viper.GetString("client_secret") == "" || viper.GetString("redirect_uri") == "" {
		m.configInputs[0].SetValue(viper.GetString("client_id"))
		m.configInputs[1].SetValue(viper.GetString("client_secret"))
		m.configInputs[2].SetValue(viper.GetString("redirect_uri"))
		cmds = append(cmds, m.changeView(models.ConfigView))
	} else if viper.GetString("access_token") == "" {
		cmds = append(cmds, m.changeView(models.AuthView))
	} else { // Load projects
		m.client = client.NewClient(viper.GetString("access_token"))
		cmds = append(cmds, m.changeView(models.ProjectListView))
	}

	cmds = append(cmds, textinput.Blink, m.spinner.Tick)
	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Message handling
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Force a re-render by returning the model
		return m, nil

	case tea.KeyMsg:
		// Handle manual resize trigger for Windows
		// FIXME:
		if msg.String() == "ctrl+r" {
			// Get current terminal size
			width, height, err := term.GetSize(int(os.Stdout.Fd()))
			if err == nil {
				m.width = width
				m.height = height
			}
			return m, tea.Sequence(tea.ClearScreen, tea.EnterAltScreen)
		}
		// Typical key handling
		cmds = append(cmds, m.handleKey(msg.String()))

	case projectsLoadedMsg:
		m.state.Error = ""   // Clear any previous error
		m.state.Message = "" // Clear any previous message
		m.state.Projects = []models.Project(msg)
		items := make([]any, len(m.state.Projects))
		for i, project := range m.state.Projects {
			items[i] = project
		}
		m.state.CurrentItems = items

	case tasksLoadedMsg:
		m.state.Error = ""   // Clear any previous error
		m.state.Message = "" // Clear any previous message
		m.state.Tasks = []models.Task(msg)
		items := make([]any, len(m.state.Tasks))
		for i, task := range m.state.Tasks {
			items[i] = task
		}
		m.state.CurrentItems = items

	// case taskCreatedMsg:
	// 	m.state.Message = "任务创建成功"
	// 	m.state.CurrentView = models.TaskListView
	// 	return m, m.loadTasks(m.state.CurrentProject.ID)
	// case projectCreatedMsg:
	// 	m.state.Message = "项目创建成功"
	// 	m.state.CurrentView = models.ProjectListView
	// 	return m, m.loadProjects()

	case configSavedMsg:
		m.state.Message = "配置已保存！"
		return m, m.changeView(models.AuthView)

	case tokenExchangedMsg:
		m.state.Message = "认证成功！"
		m.client = client.NewClient(viper.GetString("access_token"))
		return m, m.changeView(models.ProjectListView)
	}

	var cmd tea.Cmd

	// Update spinner
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	// Update items
	for idx := range m.state.CurrentItems {
		switch v := m.state.CurrentItems[idx].(type) {
		case textinput.Model:
			updatedInput, itemCmd := v.Update(msg)
			cmds = append(cmds, itemCmd)

			if idx == m.state.SelectedIndex {
				// Focus the selected input
				updatedInput.Focus()
			} else {
				// Blur other inputs
				updatedInput.Blur()
			}
			m.state.CurrentItems[idx] = updatedInput
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	statusBar := m.renderStatusBar()

	var content string
	switch m.state.CurrentView {
	case models.ConfigView:
		content = m.renderConfigForm()
	case models.AuthView:
		content = m.renderAuthForm()
	case models.ProjectListView:
		content = m.renderProjectList()
	case models.TaskListView:
		content = m.renderTaskList()
	}

	message := m.renderMessage()
	help := m.renderHelp()

	// Calculate content height to fill remaining space
	contentHeight := m.height - 1 - 1 - 2 // Subtract status bar, message line, and help bar

	// Ensure content has minimum height
	if contentHeight < 3 {
		contentHeight = 3
	}

	// Apply height to content if needed
	if content != "" && contentHeight > 0 {
		content = lipgloss.NewStyle().Height(contentHeight).Render(content)
	}

	var parts []string
	parts = append(parts, statusBar)
	if content != "" {
		parts = append(parts, content)

	}
	if message == "" {
		message = "\n" // Empty line to reserve space
	}
	parts = append(parts, message)
	parts = append(parts, help)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
