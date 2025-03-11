package models

// CreateNewTasksRequest represents the request body for creating new tasks.
type CreateNewTasksRequest struct {
	Tasks []NewTask `json:"tasks"`
}

// NewTask represents a new task for the client to create new tasks.
type NewTask struct {
	Name   string `json:"name" validate:"required" example:"Task 1"`
	Status int    `json:"status" validate:"required" example:"1" enums:"0,1"`
}

// Validate validates the task name and status and returns an error if the name is empty or the status is invalid
func (nt *NewTask) Validate() error {
	if err := nt.ValidateName(); err != nil {
		return err
	}

	if err := nt.ValidateStatus(); err != nil {
		return err
	}

	return nil
}

// ValidateName validates the task name and returns an error if the name is empty
func (nt *NewTask) ValidateName() error {
	if nt.Name == "" {
		return ErrTaskNameEmpty
	}

	return nil
}

// ValidateStatus validates the task status and returns an error if the status is invalid (not 0 or 1)
func (nt *NewTask) ValidateStatus() error {
	switch nt.Status {
	case 0, 1:
		return nil
	}

	return ErrInvalidStatus
}

// UpdateTaskRequest represents the request body for updating an existing task.
type UpdateTaskRequest struct {
	Name   *string `json:"name,omitempty"`
	Status *int    `json:"status,omitempty" enums:"0,1"`
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
