package models

// CreateNewTasksRequest represents the request body for creating new tasks.
type CreateNewTasksRequest struct {
	Tasks []Task `json:"tasks"`
}

// UpdateTaskRequest represents the request body for updating an existing task.
type UpdateTaskRequest struct {
	Name   *string `json:"name,omitempty"`
	Status *int    `json:"status,omitempty"`
}

// IsNoChanges checks if the UpdateTaskRequest contains no changes.
func (utr *UpdateTaskRequest) IsNoChanges() bool {
	return utr.Name == nil && utr.Status == nil
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
