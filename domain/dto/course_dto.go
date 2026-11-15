package dto

import "github.com/valentinusdelvin/studytrack-be/domain/entity"

type CreateCourseRequest struct {
	Name     string `json:"name" validate:"required,max=255"`
	Code     string `json:"kode" validate:"max=64"`
	Lecturer string `json:"dosen" validate:"max=255"`
}

type UpdateCourseRequest struct {
	Name     *string `json:"name" validate:"omitempty,max=255"`
	Code     *string `json:"kode" validate:"omitempty,max=64"`
	Lecturer *string `json:"dosen" validate:"omitempty,max=255"`
}

type CourseResponse struct {
	ID       uint64 `json:"id"`
	UserID   uint64 `json:"user_id"`
	Name     string `json:"name"`
	Code     string `json:"kode"`
	Lecturer string `json:"dosen"`
}

type CoursesResponse struct {
	Courses []CourseResponse `json:"courses"`
}

type CourseEnvelope struct {
	Course CourseResponse `json:"course"`
}

func FormatToCourseResponse(course entity.Course) CourseResponse {
	return CourseResponse{
		ID: course.ID, UserID: course.UserID, Name: course.Name,
		Code: course.Code, Lecturer: course.Lecturer,
	}
}
