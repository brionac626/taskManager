package models

import (
	"errors"
	"sort"

	"github.com/rs/xid"
)

var (
	// ErrTaskNameEmpty represents an error when the task name is empty
	ErrTaskNameEmpty = errors.New("task name is empty")
	// ErrInvalidStatus represents an error when the task status is invalid (not 0 or 1)
	ErrInvalidStatus = errors.New("invalid status")
)

// Task represents a task
type Task struct {
	ID     string `json:"id"`     // task id
	Name   string `json:"name"`   // task name
	Status int    `json:"status"` // 0 represents an incomplete task, while 1 represents a completed task
}

// NewTaskID generates a new task id
func (t *Task) NewTaskID() {
	t.ID = xid.New().String()
}

// Validate validates the task name and status and returns an error if the name is empty or the status is invalid
func (t *Task) Validate() error {
	if err := t.ValidateName(); err != nil {
		return err
	}

	if err := t.ValidateStatus(); err != nil {
		return err
	}

	return nil
}

// ValidateName validates the task name and returns an error if the name is empty
func (t *Task) ValidateName() error {
	if t.Name == "" {
		return ErrTaskNameEmpty
	}

	return nil
}

// ValidateStatus validates the task status and returns an error if the status is invalid (not 0 or 1)
func (t *Task) ValidateStatus() error {
	switch t.Status {
	case 0, 1:
		return nil
	}

	return ErrInvalidStatus
}

// TasksByID attaches the methods of sort.Interface to []Task, sorting by ID
type TasksByID []Task

func (t TasksByID) Len() int           { return len(t) }
func (t TasksByID) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TasksByID) Less(i, j int) bool { return t[i].ID < t[j].ID }

// SortTasksByID sorts the tasks by ID in ascending order
func SortTasksByID(tasks []Task) {
	sort.Sort(TasksByID(tasks))
}

// SortTasksByIDReverse sorts the tasks by ID in descending order
func SortTasksByIDReverse(tasks []Task) {
	sort.Sort(sort.Reverse(TasksByID(tasks)))
}
