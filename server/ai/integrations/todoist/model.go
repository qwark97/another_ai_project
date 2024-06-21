package todoist

import "time"

type GetActiveTasksParams struct {
	Date string `json:"date"`
}

type Task struct {
	ID        string    `json:"id,omitempty"`
	ProjectID string    `json:"project_id,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DueString string    `json:"due_string,omitempty"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
