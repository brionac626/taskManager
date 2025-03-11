package repository

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/brionac626/taskManager/models"
	"github.com/stretchr/testify/assert"
)

var testRepo = &taskRepo{}

func Test_taskRepo_GetTasks(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	expectTasks := []models.Task{
		{ID: "task1", Name: "Task 1", Status: 0},
		{ID: "task2", Name: "Task 2", Status: 1},
	}
	expectNoTasks := make([]models.Task, 0)
	canceledCtx, cancel := context.WithCancel(context.Background())

	tests := []struct {
		name           string
		mockSetup      func()
		want           []models.Task
		args           args
		wantErr        bool
		wantErrContent error
	}{
		{
			name: "get tasks",
			mockSetup: func() {
				for _, task := range expectTasks {
					manager.Store(task.ID, task)
				}
			},
			want:    expectTasks,
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
		{
			name: "get tasks failed (type error)",
			mockSetup: func() {
				manager.Store("task3", nil)
			},
			want:           expectNoTasks,
			args:           args{ctx: context.Background()},
			wantErr:        true,
			wantErrContent: ErrGetTasksFailed,
		},
		{
			name: "get tasks with context cancellation",
			mockSetup: func() {
				cancel()
			},
			want:           expectNoTasks,
			args:           args{ctx: canceledCtx},
			wantErr:        true,
			wantErrContent: context.Canceled,
		},
		{
			name: "no tasks",
			mockSetup: func() {
				// reset the map
				manager = sync.Map{}
			},
			want:    expectNoTasks,
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			got, err := testRepo.GetTasks(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskRepo.GetTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_taskRepo_CreateTasks(t *testing.T) {
	type args struct {
		ctx   context.Context
		tasks []models.Task
	}

	newTasks := []models.Task{
		{Name: "Task 1", Status: 0},
		{Name: "Task 2", Status: 1},
	}

	canceledCtx, cancel := context.WithCancel(context.Background())

	tests := []struct {
		name           string
		mockSetup      func()
		args           args
		wantErr        bool
		wantErrContent error
	}{
		{
			name: "create tasks",
			args: args{
				ctx:   context.Background(),
				tasks: newTasks,
			},
			wantErr: false,
		},
		{
			name: "create tasks with context cancellation",
			mockSetup: func() {
				cancel()
			},
			args: args{
				ctx:   canceledCtx,
				tasks: newTasks,
			},
			wantErr:        true,
			wantErrContent: context.Canceled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			if err := testRepo.CreateTasks(tt.args.ctx, tt.args.tasks); ((err != nil) != tt.wantErr) && errors.Is(err, tt.wantErrContent) {
				t.Errorf("taskRepo.CreateTasks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_taskRepo_DeleteTask(t *testing.T) {
	type args struct {
		ctx    context.Context
		taskID string
	}

	deleteTaskID := "task1"
	task := models.Task{ID: "task1", Name: "Task 1", Status: 0}
	canceledCtx, cancel := context.WithCancel(context.Background())

	tests := []struct {
		name           string
		mockSetup      func()
		args           args
		wantErr        bool
		wantErrContent error
	}{
		{
			name: "delete task",
			mockSetup: func() {
				manager.Store(deleteTaskID, task)
			},
			args: args{
				ctx:    context.Background(),
				taskID: deleteTaskID,
			},
			wantErr: false,
		},
		{
			name: "delete task with non-existing ID",
			mockSetup: func() {
				manager.Store(deleteTaskID, task)
			},
			args: args{
				ctx:    context.Background(),
				taskID: "non-existing-task-id",
			},
			wantErr:        true,
			wantErrContent: ErrTaskNotFound,
		},
		{
			name: "delete task with context cancellation",
			mockSetup: func() {
				cancel()
			},
			args: args{
				ctx:    canceledCtx,
				taskID: deleteTaskID,
			},
			wantErr:        true,
			wantErrContent: context.Canceled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			if err := testRepo.DeleteTask(tt.args.ctx, tt.args.taskID); ((err != nil) != tt.wantErr) && !errors.Is(err, tt.wantErrContent) {
				t.Errorf("taskRepo.DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_taskRepo_UpdateTask(t *testing.T) {
	type args struct {
		ctx    context.Context
		taskID string
		name   *string
		status *int
	}

	targetTask := models.Task{ID: "task1", Name: "Task 1", Status: 0}
	expectedUpdatedTaskName := "Updated Task 1"
	expectedUpdatedTaskStatus := 1
	expectedUpdatedTask := models.Task{ID: targetTask.ID, Name: expectedUpdatedTaskName, Status: expectedUpdatedTaskStatus}
	expectedUpdatedTaskOnlyName := models.Task{ID: targetTask.ID, Name: expectedUpdatedTaskName, Status: targetTask.Status}
	expectedUpdatedTaskOnlyStatus := models.Task{ID: targetTask.ID, Name: targetTask.Name, Status: expectedUpdatedTaskStatus}
	canceledCtx, cancel := context.WithCancel(context.Background())

	tests := []struct {
		name           string
		mockSetup      func()
		args           args
		want           *models.Task
		wantErr        bool
		wantErrContent error
	}{
		{
			name: "update task with name and status",
			mockSetup: func() {
				manager = sync.Map{}
				manager.Store(targetTask.ID, targetTask)
			},
			args: args{
				ctx:    context.Background(),
				taskID: targetTask.ID,
				name:   &expectedUpdatedTaskName,
				status: &expectedUpdatedTaskStatus,
			},
			want:    &expectedUpdatedTask,
			wantErr: false,
		},
		{
			name: "update task with only name",
			mockSetup: func() {
				manager = sync.Map{}
				manager.Store(targetTask.ID, targetTask)
			},
			args: args{
				ctx:    context.Background(),
				taskID: targetTask.ID,
				name:   &expectedUpdatedTaskName,
				status: nil,
			},
			want:    &expectedUpdatedTaskOnlyName,
			wantErr: false,
		},
		{
			name: "update task with only status",
			mockSetup: func() {
				manager = sync.Map{}
				manager.Store(targetTask.ID, targetTask)
			},
			args: args{
				ctx:    context.Background(),
				taskID: targetTask.ID,
				name:   nil,
				status: &expectedUpdatedTaskStatus,
			},
			want:    &expectedUpdatedTaskOnlyStatus,
			wantErr: false,
		},
		{
			name: "update task with non-existing ID", mockSetup: func() {
				manager = sync.Map{}
			},
			args: args{
				ctx:    context.Background(),
				taskID: "non-existing-task-id",
				name:   &expectedUpdatedTaskName,
				status: &expectedUpdatedTaskStatus,
			},
			wantErr:        true,
			wantErrContent: ErrTaskNotFound,
		},
		{
			name: "update task with context cancellation",
			mockSetup: func() {
				cancel()
			},
			args: args{
				ctx:    canceledCtx,
				taskID: targetTask.ID,
				name:   &expectedUpdatedTaskName,
				status: &expectedUpdatedTaskStatus,
			},
			wantErr:        true,
			wantErrContent: context.Canceled,
		},
		{
			name: "update task with wrong type",
			mockSetup: func() {
				manager = sync.Map{}
				manager.Store(targetTask.ID, nil)
			},
			args: args{
				ctx:    context.Background(),
				taskID: targetTask.ID,
				name:   &expectedUpdatedTaskName,
				status: &expectedUpdatedTaskStatus,
			},
			wantErr:        true,
			wantErrContent: ErrTaskType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			if err := testRepo.UpdateTask(tt.args.ctx, tt.args.taskID, tt.args.name, tt.args.status); ((err != nil) != tt.wantErr) && !errors.Is(err, tt.wantErrContent) {
				t.Errorf("taskRepo.UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				task, _ := manager.Load(targetTask.ID)
				assert.Equal(t, *tt.want, task.(models.Task))
			}
		})
	}
}
