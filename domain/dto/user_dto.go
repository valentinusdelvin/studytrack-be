package dto

import "github.com/valentinusdelvin/studytrack-be/domain/entity"

type UserResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"nama"`
	Email string `json:"email"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

func FormatToUserResponse(user entity.User) UserResponse {
	return UserResponse{ID: user.ID, Name: user.Name, Email: user.Email}
}
