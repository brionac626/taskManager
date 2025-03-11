package taskmanager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brionac626/taskManager/internal/repository"
	mocks "github.com/brionac626/taskManager/internal/repository/mocks"
	"github.com/brionac626/taskManager/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_GetTasks(t *testing.T) {
	mockTM := mocks.NewMockTaskManager(gomock.NewController(t))
	handler := &Handler{repo: mockTM}

	e := echo.New()
	e.GET("/tasks", handler.GetTasks)

	expectedTasks := []models.Task{
		{ID: "1", Name: "Task 1", Status: 0},
		{ID: "2", Name: "Task 2", Status: 1},
	}

	type args struct {
		req *http.Request
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name               string
		mockSetup          func()
		args               args
		expectedStatusCode int
		expectedResponse   any
		wantErr            bool
	}{
		{
			name: "get tasks",
			mockSetup: func() {
				mockTM.EXPECT().GetTasks(context.Background()).Return(expectedTasks, nil)
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/tasks", nil),
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   expectedTasks,
			wantErr:            false,
		},
		{
			name: "get tasks with internal error",
			mockSetup: func() {
				mockTM.EXPECT().GetTasks(context.Background()).Return(nil, repository.ErrGetTasksFailed)
			},
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/tasks", nil),
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   models.ErrorResponse{Code: http.StatusInternalServerError},
			wantErr:            false,
		},
		{
			name: "get tasks with wrong path",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/tasks123", nil),
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
			wantErr:            true,
		},
		{
			name: "get tasks with wrong method",
			args: args{
				req: httptest.NewRequest(http.MethodPatch, "/tasks", nil),
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			e.ServeHTTP(tt.args.rec, tt.args.req)

			assert.Equal(t, tt.expectedStatusCode, tt.args.rec.Result().StatusCode)

			body, err := io.ReadAll(tt.args.rec.Body)
			assert.NoError(t, err)

			if tt.wantErr {
				t.Log(string(body))
				return
			}

			if tt.args.rec.Result().StatusCode != http.StatusOK {
				if tt.expectedResponse != nil {
					var resp models.ErrorResponse
					err := json.Unmarshal(body, &resp)
					assert.NoError(t, err)
					assert.Equal(t, tt.args.rec.Result().StatusCode, resp.Code)
					return
				}
			} else {
				var tasks []models.Task
				err := json.Unmarshal(body, &tasks)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse.([]models.Task), tasks)
			}
		})
	}
}

func TestHandler_CreateTasks(t *testing.T) {
	mockTM := mocks.NewMockTaskManager(gomock.NewController(t))
	handler := &Handler{repo: mockTM}

	e := echo.New()
	e.POST("/tasks", handler.CreateTasks)

	tasks := models.CreateNewTasksRequest{
		Tasks: []models.Task{
			{ID: "1", Name: "Task 1", Status: 0},
			{ID: "2", Name: "Task 2", Status: 1},
		},
	}
	reqBody, err := json.Marshal(tasks)
	assert.NoError(t, err)
	invalidTasks := []models.Task{
		{ID: "", Name: "Task 1", Status: 0},
		{ID: "", Name: "Task 2", Status: 1},
	}
	invalidReqBody, err := json.Marshal(invalidTasks)
	assert.NoError(t, err)
	invalidTaskName := models.CreateNewTasksRequest{
		Tasks: []models.Task{
			{ID: "1", Name: "", Status: 0},
			{ID: "2", Name: "", Status: 0},
		},
	}
	invalidNameReqBody, err := json.Marshal(invalidTaskName)
	assert.NoError(t, err)
	invalidTaskStatus := models.CreateNewTasksRequest{
		Tasks: []models.Task{
			{ID: "1", Name: "Task 1", Status: 3},
			{ID: "2", Name: "Task 2", Status: -1},
		},
	}
	invalidStatusReqBody, err := json.Marshal(invalidTaskStatus)
	assert.NoError(t, err)

	type args struct {
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name               string
		mockSetup          func() *http.Request
		args               args
		expectedStatusCode int
		expectedResponse   any
		wantErr            bool
	}{
		{
			name: "create tasks",
			mockSetup: func() *http.Request {
				mockTM.EXPECT().CreateTasks(context.Background(), tasks.Tasks).Return(nil)

				req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusCreated,
			wantErr:            false,
		},
		{
			name: "create tasks with invalid request body",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString("invalid_json"))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.ErrorResponse{Code: http.StatusBadRequest},
			wantErr:            false,
		},
		{
			name: "create tasks with no content type header",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(reqBody))

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.ErrorResponse{Code: http.StatusBadRequest},
			wantErr:            false,
		},
		{
			name: "create tasks with empty request body",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.ErrorResponse{Code: http.StatusBadRequest},
			wantErr:            false,
		},
		{
			name: "create tasks with invalid tasks",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(invalidReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.ErrorResponse{Code: http.StatusBadRequest},
			wantErr:            false,
		},
		{
			name: "create tasks with invalid tasks name",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(invalidNameReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.ErrorResponse{Code: http.StatusBadRequest},
			wantErr:            false,
		},
		{
			name: "create tasks with invalid tasks status",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(invalidStatusReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.ErrorResponse{Code: http.StatusBadRequest},
			wantErr:            false,
		},
		{
			name: "create tasks with invalid tasks",
			mockSetup: func() *http.Request {
				mockTM.EXPECT().CreateTasks(context.Background(), tasks.Tasks).Return(repository.ErrTaskID)

				req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   models.ErrorResponse{Code: http.StatusInternalServerError},
			wantErr:            false,
		},
		{
			name: "create tasks with wrong path",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/tasks1234", bytes.NewBuffer(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
			wantErr:            true,
		},
		{
			name: "create tasks with wrong method",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPatch, "/tasks", bytes.NewBuffer(reqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.mockSetup()

			e.ServeHTTP(tt.args.rec, req)

			assert.Equal(t, tt.expectedStatusCode, tt.args.rec.Result().StatusCode)

			body, err := io.ReadAll(tt.args.rec.Body)
			assert.NoError(t, err)

			if tt.wantErr {
				t.Log(string(body))
				return
			}

			if tt.args.rec.Result().StatusCode == http.StatusCreated {
				return
			}

			if tt.expectedResponse != nil {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.rec.Result().StatusCode, resp.Code)
			}
		})
	}
}

func TestHandler_UpdateTask(t *testing.T) {
	mockTM := mocks.NewMockTaskManager(gomock.NewController(t))
	handler := &Handler{repo: mockTM}

	e := echo.New()
	e.PUT("/tasks/:id", handler.UpdateTask)

	taskID := "1"
	notFoundTaskID := "not_found"
	name := "Updated Task"
	status := 1
	invalidName := ""
	invalidStatus := -1
	updateNameAndStatus := models.UpdateTaskRequest{Name: &name, Status: &status}
	nameAndStatusReqBody, err := json.Marshal(updateNameAndStatus)
	assert.NoError(t, err)
	updateNameOnly := models.UpdateTaskRequest{Name: &name}
	nameOnlyReqBody, err := json.Marshal(updateNameOnly)
	assert.NoError(t, err)
	updateStatusOnly := models.UpdateTaskRequest{Status: &status}
	statusOnlyReqBody, err := json.Marshal(updateStatusOnly)
	assert.NoError(t, err)
	updateInvalidName := models.UpdateTaskRequest{Name: &invalidName}
	updateInvalidNameReqBody, err := json.Marshal(updateInvalidName)
	assert.NoError(t, err)
	updateInvalidStatus := models.UpdateTaskRequest{Status: &invalidStatus}
	updateInvalidStatusReqBody, err := json.Marshal(updateInvalidStatus)
	assert.NoError(t, err)
	updateInvalidNameAndStatus := models.UpdateTaskRequest{Name: &invalidName, Status: &invalidStatus}
	updateInvalidNameAndStatusReqBody, err := json.Marshal(updateInvalidNameAndStatus)
	assert.NoError(t, err)

	type args struct {
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name               string
		mockSetup          func() *http.Request
		args               args
		expectedStatusCode int
		expectedResponse   any
		wantErr            bool
	}{
		{
			name: "update task",
			mockSetup: func() *http.Request {
				mockTM.EXPECT().UpdateTask(context.Background(), taskID, &name, &status).Return(nil)

				req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%s", taskID), bytes.NewBuffer(nameAndStatusReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "update task name only",
			mockSetup: func() *http.Request {
				mockTM.EXPECT().UpdateTask(context.Background(), taskID, &name, nil).Return(nil)

				req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%s", taskID), bytes.NewBuffer(nameOnlyReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "update task status only",
			mockSetup: func() *http.Request {
				mockTM.EXPECT().UpdateTask(context.Background(), taskID, nil, &status).Return(nil)

				req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%s", taskID), bytes.NewBuffer(statusOnlyReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "update task no changes",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%s", taskID), nil)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "update task invalid request body",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%s", taskID), bytes.NewBufferString("invalid_json"))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.ErrorResponse{Code: http.StatusBadRequest},
			wantErr:            false,
		},
		{
			name: "update task invalid task name",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%s", taskID), bytes.NewBuffer(updateInvalidNameReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.ErrorResponse{Code: http.StatusBadRequest},
			wantErr:            false,
		},
		{
			name: "update task invalid task status",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%s", taskID), bytes.NewBuffer(updateInvalidStatusReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.ErrorResponse{Code: http.StatusBadRequest},
			wantErr:            false,
		},
		{
			name: "update task invalid task name and status",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%s", taskID), bytes.NewBuffer(updateInvalidNameAndStatusReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   models.ErrorResponse{Code: http.StatusBadRequest},
			wantErr:            false,
		},
		{
			name: "update task not found",
			mockSetup: func() *http.Request {
				mockTM.EXPECT().UpdateTask(context.Background(), notFoundTaskID, &name, &status).Return(repository.ErrTaskNotFound)

				req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks/%s", notFoundTaskID), bytes.NewBuffer(nameAndStatusReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   models.ErrorResponse{Code: http.StatusInternalServerError},
			wantErr:            false,
		},
		{
			name: "update tasks with wrong path",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/tasks1234/%s", taskID), bytes.NewBuffer(nameAndStatusReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
			wantErr:            true,
		},
		{
			name: "update tasks with wrong method",
			mockSetup: func() *http.Request {
				req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/tasks/%s", taskID), bytes.NewBuffer(nameAndStatusReqBody))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

				return req
			},
			args: args{
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.mockSetup()

			e.ServeHTTP(tt.args.rec, req)

			assert.Equal(t, tt.expectedStatusCode, tt.args.rec.Result().StatusCode)

			body, err := io.ReadAll(tt.args.rec.Body)
			assert.NoError(t, err)

			if tt.wantErr {
				t.Log(string(body))
				return
			}

			if tt.args.rec.Result().StatusCode == http.StatusCreated {
				return
			}

			if tt.expectedResponse != nil {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.rec.Result().StatusCode, resp.Code)
			}
		})
	}
}

func TestHandler_DeleteTask(t *testing.T) {
	mockTM := mocks.NewMockTaskManager(gomock.NewController(t))
	handler := &Handler{repo: mockTM}

	e := echo.New()
	e.DELETE("/tasks/:id", handler.DeleteTask)

	taskID := "1"

	type args struct {
		req *http.Request
		rec *httptest.ResponseRecorder
	}
	tests := []struct {
		name               string
		mockSetup          func()
		args               args
		expectedStatusCode int
		expectedResponse   any
		wantErr            bool
	}{
		{
			name: "delete task",
			mockSetup: func() {
				mockTM.EXPECT().DeleteTask(context.Background(), taskID).Return(nil)
			},
			args: args{
				req: httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%s", taskID), nil),
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   nil,
			wantErr:            false,
		},
		{
			name: "delete task not found",
			mockSetup: func() {
				mockTM.EXPECT().DeleteTask(context.Background(), taskID).Return(repository.ErrTaskNotFound)
			},
			args: args{
				req: httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks/%s", taskID), nil),
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   models.ErrorResponse{Code: http.StatusInternalServerError},
			wantErr:            false,
		},
		{
			name: "update tasks with wrong path",
			args: args{
				req: httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/tasks1234/%s", taskID), nil),
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusNotFound,
			wantErr:            true,
		},
		{
			name: "update tasks with wrong method",
			args: args{
				req: httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/tasks/%s", taskID), nil),
				rec: httptest.NewRecorder(),
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			e.ServeHTTP(tt.args.rec, tt.args.req)

			assert.Equal(t, tt.expectedStatusCode, tt.args.rec.Result().StatusCode)

			body, err := io.ReadAll(tt.args.rec.Body)
			assert.NoError(t, err)

			if tt.wantErr {
				t.Log(string(body))
				return
			}

			if tt.args.rec.Result().StatusCode == http.StatusCreated {
				return
			}

			if tt.expectedResponse != nil {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.args.rec.Result().StatusCode, resp.Code)
			}
		})
	}
}
