package tui

import (
	"ticktick-tui/internal/core"
	"ticktick-tui/internal/models"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleKey(key string) tea.Cmd {
	switch key {
	case "ctrl+c":
		return tea.Quit
	case "up":
		m.moveSelection(-1)
	case "down":
		m.moveSelection(1)
		// case "a":
		// case "d":
		// 	return m.handleDelete()
	case "enter":
		return m.handleComplete()
	}
	return nil
}

func (m *Model) handleDelete() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) handleComplete() tea.Cmd {
	m.state.Error = ""     // Clear any previous error
	m.state.Message = ""   // Clear any previous message
	m.state.Loading = true // Set loading state
	switch m.state.CurrentView {
	case models.ConfigView:
		return func() tea.Msg {
			defer func() {
				m.state.Loading = false
			}()

			// CurrentItems -> configInputs
			for i, input := range m.state.CurrentItems {
				if textInput, ok := input.(textinput.Model); ok {
					m.configInputs[i].SetValue(textInput.Value())
				}
			}
			// Validate inputs
			for _, input := range m.configInputs {
				if input.Value() == "" {
					m.state.Error = "All fields must be filled out."
					return nil
				}
			}

			// Save configuration
			core.SaveConfig("client_id", m.configInputs[0].Value())
			core.SaveConfig("client_secret", m.configInputs[1].Value())
			core.SaveConfig("redirect_uri", m.configInputs[2].Value())

			return configSavedMsg{}
		}

	case models.AuthView:
		return func() tea.Msg {
			defer func() {
				m.state.Loading = false
			}()

			// CurrentItems -> authInputs
			for i, input := range m.state.CurrentItems {
				if textInput, ok := input.(textinput.Model); ok {
					m.authInputs[i].SetValue(textInput.Value())
				}
			}
			// Validate inputs
			for _, input := range m.authInputs {
				if input.Value() == "" {
					m.state.Error = "Authorization code must be filled out."
					return nil
				}
			}

			// Exchange authorization code for access token
			accessToken, err := core.GetToken(m.authInputs[0].Value())
			if err != nil {
				m.state.Error = "Failed to exchange authorization code: " + err.Error()
				return nil
			}

			// Save authorization code
			core.SaveConfig("access_token", accessToken.AccessToken)

			return tokenExchangedMsg{}
		}

	case models.ProjectListView:
		m.state.Loading = false
		return m.changeView(models.TaskListView)
	}
	return nil
}

func (m *Model) moveSelection(direction int) {
	maxIndex := len(m.state.CurrentItems) - 1

	// Wrap around
	m.state.SelectedIndex += direction
	if m.state.SelectedIndex < 0 {
		m.state.SelectedIndex = maxIndex
	} else if m.state.SelectedIndex > maxIndex {
		m.state.SelectedIndex = 0
	}
}

func (m *Model) generateAuthURL() {
	m.state.AuthURL = core.GetAuthURL()
}

func (m *Model) resetForm() {
	// m.projectInput.SetValue("")
	// m.projectInput.Focus()

	for _, input := range m.configInputs {
		input.SetValue("")
	}
	for _, input := range m.authInputs {
		input.SetValue("")
	}
}

func (m *Model) loadProjects() tea.Cmd {
	return func() tea.Msg {
		m.state.Loading = true
		defer func() {
			m.state.Loading = false
		}()
		projects, err := core.GetProjects()
		if err != nil {
			m.state.Error = "Failed to load projects: " + err.Error()
			return nil
		}

		// Sort projects first by GroupID, then by SortOrder within each group
		for i := 0; i < len(projects); i++ {
			for j := i + 1; j < len(projects); j++ {
				// First compare GroupID
				if projects[i].GroupID > projects[j].GroupID {
					projects[i], projects[j] = projects[j], projects[i]
				} else if projects[i].GroupID == projects[j].GroupID {
					// If GroupID is the same, compare SortOrder
					if projects[i].SortOrder > projects[j].SortOrder {
						projects[i], projects[j] = projects[j], projects[i]
					}
				}
			}
		}
		return projectsLoadedMsg(projects)
	}
}

func (m *Model) loadTasks() tea.Cmd {
	return func() tea.Msg {
		m.state.Loading = true
		defer func() {
			m.state.Loading = false
		}()
		if m.state.CurrentProject == nil {
			m.state.Error = "No project selected."
			return nil
		}
		tasks, err := core.GetTasks(m.state.CurrentProject.ID)
		if err != nil {
			m.state.Error = "Failed to load tasks: " + err.Error()
			return nil
		}
		return tasksLoadedMsg(tasks)
	}
}

func (m *Model) changeView(view models.ViewState) tea.Cmd {
	m.state.CurrentView = view
	defer func() {
		m.state.SelectedIndex = 0
	}()
	m.state.CurrentItems = []any{}

	switch view {
	case models.ConfigView:
		items := make([]any, len(m.configInputs))
		for i, input := range m.configInputs {
			items[i] = input
		}
		m.state.CurrentItems = items

	case models.AuthView:
		m.generateAuthURL()
		// Copy to clipboard
		clipboard.WriteAll(m.state.AuthURL)
		m.resetForm()
		items := make([]any, len(m.authInputs))
		for i, input := range m.authInputs {
			items[i] = input
		}
		m.state.CurrentItems = items

	case models.ProjectListView:
		m.resetForm()
		m.state.CurrentItems = []any{}
		return m.loadProjects()

	case models.TaskListView:
		m.resetForm()
		if len(m.state.Projects) == 0 {
			m.state.Error = "No projects available."
			return nil
		}
		if m.state.SelectedIndex < 0 || m.state.SelectedIndex >= len(m.state.Projects) {
			m.state.Error = "Invalid project selection."
			return nil
		}
		m.state.CurrentProject = &m.state.Projects[m.state.SelectedIndex]
		m.state.CurrentItems = []any{}
		return m.loadTasks()

	case models.TaskDetailView:
		// Show task details here
	case models.CreateTaskView:
		// Show create task form here
	case models.CreateProjectView:
		// Show create project form here
	case models.DeleteConfirmView:
		// Show delete confirmation here
	}

	return nil
}
