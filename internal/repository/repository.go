package repository

import (
	"context"

	"github.com/brionac626/taskManager/models"
)

// TaskManager represents a task manager to manage tasks in the memory
type TaskManager interface {
	GetTasks(ctx context.Context) ([]models.Task, error)
	CreateTasks(ctx context.Context, tasks []models.Task) error
	UpdateTask(ctx context.Context, taskID string, name *string, status *int) error
	DeleteTask(ctx context.Context, taskID string) error
}
