package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/brionac626/taskManager/models"
	"github.com/rs/xid"
)

type taskRepo struct{}

var _ TaskManager = (*taskRepo)(nil)

// in-memory storage for tasks
var manager sync.Map

var (
	// ErrGetTasksFailed represents an error when getting tasks failed
	ErrGetTasksFailed = errors.New("failed to get tasks")
	// ErrTaskNotFound represents an error when a task is not found
	ErrTaskNotFound = errors.New("task not found")
	// ErrTaskType represents an error when the task type is invalid
	ErrTaskType = errors.New("task type error")
	// ErrTaskID represents an error when the task id is invalid (not xid)
	ErrTaskID = errors.New("invalid task id")
)

// NewRepository creates a new task manager for managing tasks in the memory
func NewRepository() TaskManager {
	return &taskRepo{}
}

// GetTasks returns all tasks from the memory
func (t *taskRepo) GetTasks(ctx context.Context) ([]models.Task, error) {
	var err error
	result := make([]models.Task, 0)

	manager.Range(func(key, value interface{}) bool {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return false
		default:
		}

		task, ok := value.(models.Task)
		if !ok {
			err = ErrGetTasksFailed
			return false
		}

		result = append(result, task)
		return true
	})

	if err != nil {
		return make([]models.Task, 0), err
	}

	models.SortTasksByID(result)

	return result, nil

}

// CreateTasks creates tasks from request
func (t *taskRepo) CreateTasks(ctx context.Context, tasks []models.Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	for _, task := range tasks {
		if _, err := xid.FromString(task.ID); err != nil {
			return ErrTaskID
		}

		manager.Store(task.ID, task)
	}

	return nil
}

// UpdateTask updates a task by task id
func (t *taskRepo) UpdateTask(ctx context.Context, taskID string, name *string, status *int) error {
	v, exists := manager.Load(taskID)
	if !exists {
		return ErrTaskNotFound
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	task, ok := v.(models.Task)
	if !ok {
		return ErrTaskType
	}

	if name != nil {
		task.Name = *name
	}

	if status != nil {
		task.Status = *status
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	manager.Swap(taskID, task)

	return nil
}

// DeleteTask deletes a task by task id
func (t *taskRepo) DeleteTask(ctx context.Context, taskID string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	_, exists := manager.LoadAndDelete(taskID)
	if !exists {
		return ErrTaskNotFound
	}

	return nil
}
