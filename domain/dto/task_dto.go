package dto

import (
	"time"

	"github.com/valentinusdelvin/studytrack-be/domain/entity"
	"github.com/valentinusdelvin/studytrack-be/domain/enums"
)

type CreateTaskRequest struct {
	CourseID    string `json:"mata_kuliah_id" validate:"required,min=1"`
	Title       string `json:"judul" validate:"required,max=255"`
	Description string `json:"deskripsi" validate:"max=10000"`
	Deadline    string `json:"deadline" validate:"required"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"judul" validate:"omitempty,max=255"`
	Description *string `json:"deskripsi" validate:"omitempty,max=10000"`
	Deadline    *string `json:"deadline" validate:"omitempty"`
	Status      *string `json:"status" validate:"omitempty,oneof=belum selesai"`
}

type TaskResponse struct {
	ID          uint64 `json:"id"`
	UserID      uint64 `json:"user_id"`
	CourseID    uint64 `json:"mata_kuliah_id"`
	Title       string `json:"judul"`
	Description string `json:"deskripsi"`
	Deadline    string `json:"deadline"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

type TaskListItem struct {
	TaskResponse
	CourseName     string `json:"mata_kuliah"`
	DeadlineStatus string `json:"status_deadline"`
}

type TasksResponse struct {
	Tasks []TaskListItem `json:"tasks"`
}
type TaskEnvelope struct {
	Task TaskResponse `json:"task"`
}

func FormatToTaskResponse(task entity.Task) TaskResponse {
	return TaskResponse{
		ID: task.ID, UserID: task.UserID, CourseID: task.CourseID,
		Title: task.Title, Description: task.Description,
		Deadline: task.Deadline.Format("2006-01-02T15:04:05"),
		Status:   string(task.Status), CreatedAt: task.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func FormatToTaskListItemResponse(task entity.Task, courseName string, deadlineStatus enums.TaskDeadlineStatus) TaskListItem {
	return TaskListItem{
		TaskResponse:   FormatToTaskResponse(task),
		CourseName:     courseName,
		DeadlineStatus: string(deadlineStatus),
	}
}
