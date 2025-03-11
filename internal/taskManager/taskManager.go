package taskmanager

import (
	"log"
	"net/http"

	"github.com/brionac626/taskManager/models"

	"github.com/labstack/echo/v4"
)

// GetTasks godoc
// @Summary      Get all tasks from the local storage.
// @Description  Get all tasks.
// @Tags         Tasks
// @Produce      json
// @Success      200  {array}  []models.Task  "tasks retrieved successfully"
// @Failure      500  {object}  models.ErrorResponse  "Filed to get tasks"
// @Router       /tasks [get]
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

// CreateTasks godoc
// @Summary      Create new tasks from the client request.
// @Description  crate new tasks.
// @Tags         Tasks
// @Accept		 json
// @Param 		 req  body  models.CreateNewTasksRequest  true  "tasks to create"
// @Success      201  "no content returned when successful"
// @Failure      400  {object}  models.ErrorResponse  "Invalid request body"
// @Failure      400  {object}  models.ErrorResponse  "No tasks provided"
// @Failure      400  {object}  models.ErrorResponse  "Invalid task fields values"
// @Failure      500  {object}  models.ErrorResponse  "Failed to create tasks"
// @Router       /tasks [post]
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

	newTasks := make([]models.Task, 0)
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
		newTasks = append(newTasks, models.Task{Name: task.Name, Status: task.Status})
	}

	if err := h.repo.CreateTasks(ctx, newTasks); err != nil {
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

// UpdateTask godoc
// @Summary      Update an existing task by task id.
// @Description  Update an existing task fields' values.
// @Tags         Tasks
// @Accept		 json
// @Param 		 id  path  string  true  "target task id"	example("9bsv0s2hf8ng030mva9g")	default("9bsv0s2hf8ng030mva9g")
// @Param 		 req  body  models.UpdateTaskRequest  true  "task fields to update"
// @Success      200  "no content returned when successful"
// @Success      200  "no content returned when no changes"
// @Failure      400  {object}  models.ErrorResponse  "Invalid request body"
// @Failure      400  {object}  models.ErrorResponse  "Invalid task fields values"
// @Failure      500  {object}  models.ErrorResponse  "Failed to update a task fields"
// @Router       /tasks/:id [put]
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

	if err := req.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			&models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		)
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

// DeleteTask godoc
// @Summary      Delete an existing task by task id.
// @Description  Delete an existing task.
// @Tags         Tasks
// @Param 		 id  path  string  true  "target task id"	example("9bsv0s2hf8ng030mva9g")	default("9bsv0s2hf8ng030mva9g")
// @Success      200  "no content returned when successful"
// @Failure      500  {object}  models.ErrorResponse  "Failed to update a task fields"
// @Router       /tasks/:id [delete]
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
