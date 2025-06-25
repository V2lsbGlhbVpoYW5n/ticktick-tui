package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ticktick-tui/internal/models"
)

const BaseURL = "https://api.ticktick.com"

type Client struct {
	accessToken string
	httpClient  *http.Client
}

// NewClient creates a new TickTick API client
func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		httpClient:  &http.Client{},
	}
}

// makeRequest performs HTTP request with proper headers
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, BaseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// GetTask retrieves a task by project ID and task ID
func (c *Client) GetTask(projectID, taskID string) (*models.Task, error) {
	endpoint := fmt.Sprintf("/open/v1/project/%s/task/%s", projectID, taskID)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var task models.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &task, nil
}

// CreateTask creates a new task
func (c *Client) CreateTask(task *models.Task) (*models.Task, error) {
	resp, err := c.makeRequest("POST", "/open/v1/task", task)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var createdTask models.Task
	if err := json.NewDecoder(resp.Body).Decode(&createdTask); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &createdTask, nil
}

// UpdateTask updates an existing task
func (c *Client) UpdateTask(taskID string, task *models.Task) (*models.Task, error) {
	endpoint := fmt.Sprintf("/open/v1/task/%s", taskID)
	resp, err := c.makeRequest("POST", endpoint, task)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var updatedTask models.Task
	if err := json.NewDecoder(resp.Body).Decode(&updatedTask); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &updatedTask, nil
}

// CompleteTask marks a task as completed
func (c *Client) CompleteTask(projectID, taskID string) error {
	endpoint := fmt.Sprintf("/open/v1/project/%s/task/%s/complete", projectID, taskID)
	resp, err := c.makeRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// DeleteTask deletes a task
func (c *Client) DeleteTask(projectID, taskID string) error {
	endpoint := fmt.Sprintf("/open/v1/project/%s/task/%s", projectID, taskID)
	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}

// GetProjects retrieves all user projects
func (c *Client) GetProjects() ([]models.Project, error) {
	resp, err := c.makeRequest("GET", "/open/v1/project", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var projects []models.Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return projects, nil
}

// GetProject retrieves a project by ID
func (c *Client) GetProject(projectID string) (*models.Project, error) {
	endpoint := fmt.Sprintf("/open/v1/project/%s", projectID)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var project models.Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &project, nil
}

// GetProjectData retrieves project with tasks and columns
func (c *Client) GetProjectData(projectID string) (*models.ProjectData, error) {
	endpoint := fmt.Sprintf("/open/v1/project/%s/data", projectID)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var projectData models.ProjectData
	if err := json.NewDecoder(resp.Body).Decode(&projectData); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &projectData, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(project *models.Project) (*models.Project, error) {
	resp, err := c.makeRequest("POST", "/open/v1/project", project)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var createdProject models.Project
	if err := json.NewDecoder(resp.Body).Decode(&createdProject); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &createdProject, nil
}

// UpdateProject updates an existing project
func (c *Client) UpdateProject(projectID string, project *models.Project) (*models.Project, error) {
	endpoint := fmt.Sprintf("/open/v1/project/%s", projectID)
	resp, err := c.makeRequest("POST", endpoint, project)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var updatedProject models.Project
	if err := json.NewDecoder(resp.Body).Decode(&updatedProject); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &updatedProject, nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(projectID string) error {
	endpoint := fmt.Sprintf("/open/v1/project/%s", projectID)
	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API error: %d", resp.StatusCode)
	}

	return nil
}
