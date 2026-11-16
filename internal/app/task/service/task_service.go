package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/valentinusdelvin/studytrack-be/domain/dto"
	"github.com/valentinusdelvin/studytrack-be/domain/entity"
	"github.com/valentinusdelvin/studytrack-be/domain/enums"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/internal/app/shared/validation"
	"github.com/valentinusdelvin/studytrack-be/internal/app/task/repository"
	"gorm.io/gorm"
)

type TaskService struct{ taskRepo *repository.TaskRepository }

func NewTaskService(taskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

func (s *TaskService) GetTasks(ctx context.Context, userID uint64) (dto.TasksResponse, error) {
	rows, err := s.taskRepo.GetTasksByUserID(ctx, userID)
	if err != nil {
		return dto.TasksResponse{}, dbError(err, "task.List")
	}
	now := time.Now()
	out := formatTaskListItems(rows, now)
	return dto.TasksResponse{Tasks: out}, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, userID, id uint64) (dto.TaskEnvelope, error) {
	row, err := s.taskRepo.GetTaskByIDAndUserID(ctx, userID, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.TaskEnvelope{}, errx.ErrNotFound
	}
	if err != nil {
		return dto.TaskEnvelope{}, dbError(err, "task.Get")
	}
	return dto.TaskEnvelope{Task: dto.FormatToTaskResponse(*row)}, nil
}

func (s *TaskService) CreateTask(ctx context.Context, userID uint64, req dto.CreateTaskRequest) (dto.TaskEnvelope, error) {
	if err := validation.Validate(req); err != nil {
		return dto.TaskEnvelope{}, err
	}
	deadline, err := parseTaskDeadline(req.Deadline)
	if err != nil {
		return dto.TaskEnvelope{}, errx.ErrBadRequest.WithMessage("Format deadline tidak valid.").WithError(err)
	}

	courseID, err := strconv.Atoi(req.CourseID)
	if err != nil {
		return dto.TaskEnvelope{}, errx.ErrBadRequest.WithMessage("Course ID tidak valid.").WithError(err)
	}

	owned, err := s.taskRepo.IsCourseOwnedByUser(ctx, userID, uint64(courseID))
	if err != nil {
		return dto.TaskEnvelope{}, dbError(err, "task.Create.checkCourse")
	}
	if !owned {
		return dto.TaskEnvelope{}, errx.ErrBadRequest.WithMessage("Mata kuliah tidak ditemukan untuk user ini.")
	}
	row := entity.Task{UserID: userID, CourseID: uint64(courseID), Title: strings.TrimSpace(req.Title), Description: req.Description, Deadline: deadline, Status: enums.TaskStatusNotDone}
	if err := s.taskRepo.CreateTask(ctx, &row); err != nil {
		return dto.TaskEnvelope{}, dbError(err, "task.Create")
	}
	return dto.TaskEnvelope{Task: dto.FormatToTaskResponse(row)}, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, userID, id uint64, req dto.UpdateTaskRequest) (dto.TaskEnvelope, error) {
	if err := validation.Validate(req); err != nil {
		return dto.TaskEnvelope{}, err
	}
	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		return dto.TaskEnvelope{}, errx.ErrBadRequest.WithMessage("Judul tugas tidak boleh kosong.")
	}
	row, err := s.taskRepo.GetTaskByIDAndUserID(ctx, userID, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.TaskEnvelope{}, errx.ErrNotFound
	}
	if err != nil {
		return dto.TaskEnvelope{}, dbError(err, "task.Update.get")
	}
	if req.Title != nil {
		row.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		row.Description = *req.Description
	}
	if req.Deadline != nil {
		row.Deadline, err = parseTaskDeadline(*req.Deadline)
		if err != nil {
			return dto.TaskEnvelope{}, errx.ErrBadRequest.WithMessage("Format deadline tidak valid.").WithError(err)
		}
	}
	if req.Status != nil {
		row.Status = enums.TaskStatus(*req.Status)
	}
	if err := s.taskRepo.UpdateTask(ctx, row); err != nil {
		return dto.TaskEnvelope{}, dbError(err, "task.Update.save")
	}
	return dto.TaskEnvelope{Task: dto.FormatToTaskResponse(*row)}, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, userID, id uint64) (dto.SuccessResponse, error) {
	if _, err := s.taskRepo.GetTaskByIDAndUserID(ctx, userID, id); errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.SuccessResponse{}, errx.ErrNotFound
	} else if err != nil {
		return dto.SuccessResponse{}, dbError(err, "task.Delete.get")
	}
	if err := s.taskRepo.DeleteTask(ctx, userID, id); err != nil {
		return dto.SuccessResponse{}, dbError(err, "task.Delete")
	}
	return dto.SuccessResponse{Success: true}, nil
}
