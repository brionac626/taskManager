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

func (utr *UpdateTaskRequest) Validate() error {
	if utr.Name != nil && *utr.Name == "" {
        return ErrTaskNameEmpty
    }

    if utr.Status != nil && (*utr.Status != 0 && *utr.Status != 1) {
        return ErrInvalidStatus
    }

    return nil
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
