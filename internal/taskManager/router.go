package taskmanager

import (
	"github.com/brionac626/taskManager/internal/repository"

	"github.com/labstack/echo/v4"
)

// Handler handles tasks using the given repository.
type Handler struct {
	repo repository.TaskManager
}

// NewRouter creates a new Echo router with task manager integration.
func NewRouter(taskManager repository.TaskManager) *echo.Echo {
	handler := &Handler{repo: taskManager}

	e := echo.New()
	e.HideBanner = true

	e.GET("/tasks", handler.GetTasks)
	e.POST("/tasks", handler.CreateTasks)
	e.PUT("/tasks/:id", handler.UpdateTask)
	e.DELETE("/tasks/:id", handler.DeleteTask)

	return e
}
