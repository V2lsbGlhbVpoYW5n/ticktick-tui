package models

// TUI models

// ViewState represents the current view state in the TUI application
type ViewState int

const (
	ConfigView ViewState = iota
	AuthView
	ProjectListView
	TaskListView
	TaskDetailView
	CreateTaskView
	CreateProjectView
	DeleteConfirmView
)

// AppState represents the application state for the TUI
type AppState struct {
	CurrentView ViewState

	Projects       []Project
	CurrentProject *Project

	Tasks       []Task
	CurrentTask *Task

	CurrentItems  []any
	SelectedIndex int

	Loading bool
	Error   string
	Message string
	AuthURL string
}
