package service

import (
	"context"
	"errors"
	"github.com/valentinusdelvin/studytrack-be/domain/dto"
	"github.com/valentinusdelvin/studytrack-be/domain/entity"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/internal/app/course/repository"
	"github.com/valentinusdelvin/studytrack-be/internal/app/shared/validation"
	"gorm.io/gorm"
	"strings"
)

type Service struct{ repo *repository.Repository }

func New(repo *repository.Repository) *Service { return &Service{repo: repo} }
func (s *Service) List(ctx context.Context, userID uint64) (dto.CoursesResponse, error) {
	rows, err := s.repo.List(ctx, userID)
	if err != nil {
		return dto.CoursesResponse{}, dbError(err, "course.List")
	}
	out := make([]dto.CourseResponse, 0, len(rows))
	for _, row := range rows {
		out = append(out, dto.FormatToCourseResponse(row))
	}
	return dto.CoursesResponse{Courses: out}, nil
}
func (s *Service) Create(ctx context.Context, userID uint64, req dto.CreateCourseRequest) (dto.CourseEnvelope, error) {
	if err := validation.Validate(req); err != nil {
		return dto.CourseEnvelope{}, err
	}
	row := entity.Course{UserID: userID, Name: strings.TrimSpace(req.Name), Code: strings.TrimSpace(req.Code), Lecturer: strings.TrimSpace(req.Lecturer)}
	if err := s.repo.Create(ctx, &row); err != nil {
		return dto.CourseEnvelope{}, dbError(err, "course.Create")
	}
	return dto.CourseEnvelope{Course: dto.FormatToCourseResponse(row)}, nil
}
func (s *Service) Update(ctx context.Context, userID, id uint64, req dto.UpdateCourseRequest) (dto.CourseEnvelope, error) {
	if err := validation.Validate(req); err != nil {
		return dto.CourseEnvelope{}, err
	}
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return dto.CourseEnvelope{}, errx.ErrBadRequest.WithMessage("Nama mata kuliah tidak boleh kosong.")
	}
	row, err := s.repo.Get(ctx, userID, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.CourseEnvelope{}, errx.ErrNotFound
	}
	if err != nil {
		return dto.CourseEnvelope{}, dbError(err, "course.Update.get")
	}
	if req.Name != nil {
		row.Name = strings.TrimSpace(*req.Name)
	}
	if req.Code != nil {
		row.Code = strings.TrimSpace(*req.Code)
	}
	if req.Lecturer != nil {
		row.Lecturer = strings.TrimSpace(*req.Lecturer)
	}
	if err := s.repo.Save(ctx, row); err != nil {
		return dto.CourseEnvelope{}, dbError(err, "course.Update.save")
	}
	return dto.CourseEnvelope{Course: dto.FormatToCourseResponse(*row)}, nil
}
func (s *Service) Delete(ctx context.Context, userID, id uint64) (dto.SuccessResponse, error) {
	if _, err := s.repo.Get(ctx, userID, id); errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.SuccessResponse{}, errx.ErrNotFound
	} else if err != nil {
		return dto.SuccessResponse{}, dbError(err, "course.Delete.get")
	}
	if err := s.repo.Delete(ctx, userID, id); err != nil {
		return dto.SuccessResponse{}, dbError(err, "course.Delete")
	}
	return dto.SuccessResponse{Success: true}, nil
}
func dbError(err error, location string) error {
	return errx.ErrInternalServer.WithError(err).WithLocation(location)
}
