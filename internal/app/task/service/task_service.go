package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/valentinusdelvin/studytrack-be/domain/dto"
	"github.com/valentinusdelvin/studytrack-be/domain/entity"
	"github.com/valentinusdelvin/studytrack-be/domain/enums"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/internal/app/shared/validation"
	"github.com/valentinusdelvin/studytrack-be/internal/app/task/repository"
	"gorm.io/gorm"
)

type Service struct{ repo *repository.Repository }

func New(repo *repository.Repository) *Service { return &Service{repo: repo} }

func (s *Service) List(ctx context.Context, userID uint64) (dto.TasksResponse, error) {
	rows, err := s.repo.List(ctx, userID)
	if err != nil {
		return dto.TasksResponse{}, dbError(err, "task.List")
	}
	out := make([]dto.TaskListItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, dto.FormatToTaskListItemResponse(row.Task, row.CourseName))
	}
	return dto.TasksResponse{Tasks: out}, nil
}

func (s *Service) Get(ctx context.Context, userID, id uint64) (dto.TaskEnvelope, error) {
	row, err := s.repo.Get(ctx, userID, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.TaskEnvelope{}, errx.ErrNotFound
	}
	if err != nil {
		return dto.TaskEnvelope{}, dbError(err, "task.Get")
	}
	return dto.TaskEnvelope{Task: dto.FormatToTaskResponse(*row)}, nil
}

func (s *Service) Create(ctx context.Context, userID uint64, req dto.CreateTaskRequest) (dto.TaskEnvelope, error) {
	if err := validation.Validate(req); err != nil {
		return dto.TaskEnvelope{}, err
	}
	deadline, err := parseDeadline(req.Deadline)
	if err != nil {
		return dto.TaskEnvelope{}, errx.ErrBadRequest.WithMessage("Format deadline tidak valid.").WithError(err)
	}
	owned, err := s.repo.IsCourseOwned(ctx, userID, req.CourseID)
	if err != nil {
		return dto.TaskEnvelope{}, dbError(err, "task.Create.checkCourse")
	}
	if !owned {
		return dto.TaskEnvelope{}, errx.ErrBadRequest.WithMessage("Mata kuliah tidak ditemukan untuk user ini.")
	}
	row := entity.Task{UserID: userID, CourseID: req.CourseID, Title: strings.TrimSpace(req.Title), Description: req.Description, Deadline: deadline, Status: enums.TaskStatusNotDone}
	if err := s.repo.Create(ctx, &row); err != nil {
		return dto.TaskEnvelope{}, dbError(err, "task.Create")
	}
	return dto.TaskEnvelope{Task: dto.FormatToTaskResponse(row)}, nil
}

func (s *Service) Update(ctx context.Context, userID, id uint64, req dto.UpdateTaskRequest) (dto.TaskEnvelope, error) {
	if err := validation.Validate(req); err != nil {
		return dto.TaskEnvelope{}, err
	}
	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		return dto.TaskEnvelope{}, errx.ErrBadRequest.WithMessage("Judul tugas tidak boleh kosong.")
	}
	row, err := s.repo.Get(ctx, userID, id)
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
		row.Deadline, err = parseDeadline(*req.Deadline)
		if err != nil {
			return dto.TaskEnvelope{}, errx.ErrBadRequest.WithMessage("Format deadline tidak valid.").WithError(err)
		}
	}
	if req.Status != nil {
		row.Status = enums.TaskStatus(*req.Status)
	}
	if err := s.repo.Save(ctx, row); err != nil {
		return dto.TaskEnvelope{}, dbError(err, "task.Update.save")
	}
	return dto.TaskEnvelope{Task: dto.FormatToTaskResponse(*row)}, nil
}

func (s *Service) Delete(ctx context.Context, userID, id uint64) (dto.SuccessResponse, error) {
	if _, err := s.repo.Get(ctx, userID, id); errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.SuccessResponse{}, errx.ErrNotFound
	} else if err != nil {
		return dto.SuccessResponse{}, dbError(err, "task.Delete.get")
	}
	if err := s.repo.Delete(ctx, userID, id); err != nil {
		return dto.SuccessResponse{}, dbError(err, "task.Delete")
	}
	return dto.SuccessResponse{Success: true}, nil
}

func parseDeadline(value string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02T15:04", "2006-01-02T15:04:05", time.RFC3339, time.RFC3339Nano} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid deadline")
}

func dbError(err error, location string) error {
	return errx.ErrInternalServer.WithError(err).WithLocation(location)
}
