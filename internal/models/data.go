package models

// Task represents a TickTick task
type Task struct {
	ID            string          `json:"id,omitempty"`
	ProjectID     string          `json:"projectId"`
	Title         string          `json:"title"`
	Content       string          `json:"content,omitempty"`
	Desc          string          `json:"desc,omitempty"`
	IsAllDay      bool            `json:"isAllDay,omitempty"`
	StartDate     *TickTickTime   `json:"startDate,omitempty"`
	DueDate       *TickTickTime   `json:"dueDate,omitempty"`
	TimeZone      string          `json:"timeZone,omitempty"`
	Reminders     []string        `json:"reminders,omitempty"`
	RepeatFlag    string          `json:"repeatFlag,omitempty"`
	Priority      TaskPriority    `json:"priority,omitempty"`
	Status        int             `json:"status,omitempty"`
	CompletedTime *TickTickTime   `json:"completedTime,omitempty"`
	SortOrder     int64           `json:"sortOrder,omitempty"`
	Items         []ChecklistItem `json:"items,omitempty"`
}

// TaskPriority represents the priority level of a task
type TaskPriority int

const (
	PriorityNone   TaskPriority = 0
	PriorityLow    TaskPriority = 1
	PriorityMedium TaskPriority = 3
	PriorityHigh   TaskPriority = 5
)

// ChecklistItem represents a subtask
type ChecklistItem struct {
	ID            string        `json:"id,omitempty"`
	Title         string        `json:"title"`
	Status        int           `json:"status,omitempty"`
	CompletedTime *TickTickTime `json:"completedTime,omitempty"`
	IsAllDay      bool          `json:"isAllDay,omitempty"`
	SortOrder     int64         `json:"sortOrder,omitempty"`
	StartDate     *TickTickTime `json:"startDate,omitempty"`
	TimeZone      string        `json:"timeZone,omitempty"`
}

// Project represents a TickTick project
type Project struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name"`
	Color      string `json:"color,omitempty"`
	SortOrder  int64  `json:"sortOrder,omitempty"`
	Closed     bool   `json:"closed,omitempty"`
	GroupID    string `json:"groupId,omitempty"`
	ViewMode   string `json:"viewMode,omitempty"`
	Permission string `json:"permission,omitempty"`
	Kind       string `json:"kind,omitempty"`
}

// Column represents a project column
type Column struct {
	ID        string `json:"id,omitempty"`
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	SortOrder int64  `json:"sortOrder,omitempty"`
}

// ProjectData represents project with tasks and columns
type ProjectData struct {
	Project Project  `json:"project"`
	Tasks   []Task   `json:"tasks"`
	Columns []Column `json:"columns"`
}

// OAuthToken represents OAuth2 token response
type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}
