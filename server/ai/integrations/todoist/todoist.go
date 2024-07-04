package todoist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultBaseURL = "https://api.todoist.com"

	tasksURI    = "/rest/v2/tasks"
	projectsURI = "/rest/v2/projects"
)

type Todoist struct {
	baseURL string
	token   string
}

func New(token string, opts ...Option) Todoist {
	t := Todoist{
		token: token,
	}
	t.applyDefaults()
	for _, opt := range opts {
		opt(&t)
	}
	return t
}

func (t *Todoist) applyDefaults() {
	t.baseURL = DefaultBaseURL
}

type Option func(*Todoist)

func WithCustomBaseURL(url string) Option {
	return func(t *Todoist) {
		t.baseURL = url
	}
}

func (t Todoist) request(ctx context.Context, method, uri string, filters map[string]string, body io.Reader) (*http.Response, error) {
	address, err := url.JoinPath(t.baseURL, uri)
	address = addFilters(address, filters)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, address, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	req.Header.Add("Content-Type", "application/json")

	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode > 299 {
		return response, fmt.Errorf(response.Status)
	}
	return response, nil
}

func addFilters(address string, filters map[string]string) string {
	filtersAmount := len(filters)
	if filtersAmount == 0 {
		return address
	}
	address += "?"
	i := 0
	for key, val := range filters {
		escapedKey := url.QueryEscape(key)
		escapedVal := url.QueryEscape(val)
		address += escapedKey + "=" + escapedVal
		i++
		if i < filtersAmount {
			address += "?"
		}
	}
	return address
}

func (t Todoist) GetActiveTasks(ctx context.Context, params GetActiveTasksParams) ([]Task, error) {
	var result []Task
	filter := map[string]string{
		"filter": t.filter(params),
	}
	response, err := t.request(ctx, http.MethodGet, tasksURI, filter, nil)
	if err != nil {
		return result, err
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (t Todoist) filter(params GetActiveTasksParams) string {
	if params.Date == "" {
		params.Date = time.Now().Format(time.DateOnly)
	}
	return params.Date + " | overdue"
}

func (t Todoist) GetProjects(ctx context.Context) ([]Project, error) {
	var result []Project
	response, err := t.request(ctx, http.MethodGet, projectsURI, nil, nil)
	if err != nil {
		return result, err
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (t Todoist) CreateTask(ctx context.Context, task Task) (Task, error) {
	var createdTask Task
	data, err := json.Marshal(task)
	if err != nil {
		return createdTask, err
	}
	body := bytes.NewReader(data)

	response, err := t.request(ctx, http.MethodPost, tasksURI, nil, body)
	if err != nil {
		return createdTask, err
	}
	err = json.NewDecoder(response.Body).Decode(&createdTask)
	if err != nil {
		return createdTask, err
	}
	return createdTask, nil
}

func (t Todoist) CloseTask(ctx context.Context, taskID string) error {
	uri, err := url.JoinPath(tasksURI, taskID, "close")
	if err != nil {
		return err
	}

	_, err = t.request(ctx, http.MethodPost, uri, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
