package taskmanager

import (
	"log"
	"net/http"

	"github.com/brionac626/taskManager/models"

	"github.com/labstack/echo/v4"
)

// GetTasks retrieves all tasks.
func (h *Handler) GetTasks(c echo.Context) error {
	ctx := c.Request().Context()

	tasks, err := h.repo.GetTasks(ctx)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, &tasks)
}

// CreateTasks creates new tasks from the client request.
func (h *Handler) CreateTasks(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.CreateNewTasksRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		)
	}

	if len(req.Tasks) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			&models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "No tasks provided",
			},
		)
	}

	for _, task := range req.Tasks {
		if err := task.Validate(); err != nil {
			log.Println("invalid err", err)
			return c.JSON(
				http.StatusBadRequest,
				&models.ErrorResponse{
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			)
		}
		task.NewTaskID()
	}

	if err := h.repo.CreateTasks(ctx, req.Tasks); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		)
	}

	return c.NoContent(http.StatusCreated)
}

// UpdateTask updates an existing task by task id.
func (h *Handler) UpdateTask(c echo.Context) error {
	ctx := c.Request().Context()

	taskID := c.Param("id")
	var req models.UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		)
	}

	if req.IsNoChanges() {
		return c.NoContent(http.StatusOK)
	}

	if err := h.repo.UpdateTask(ctx, taskID, req.Name, req.Status); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		)
	}

	return c.NoContent(http.StatusOK)
}

// DeleteTask deletes an existing task by task id.
func (h *Handler) DeleteTask(c echo.Context) error {
	ctx := c.Request().Context()

	taskID := c.Param("id")
	if err := h.repo.DeleteTask(ctx, taskID); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			&models.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		)
	}

	return c.NoContent(http.StatusOK)
}
