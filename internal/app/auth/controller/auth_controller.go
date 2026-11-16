package controller

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/valentinusdelvin/studytrack-be/domain/dto"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/internal/app/auth/service"
	"github.com/valentinusdelvin/studytrack-be/pkg/helpers/http/response"
)

type AuthController struct{ authService *service.AuthService }

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) RegisterUser(ctx fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, errx.ErrBadRequest.WithError(err))
	}
	res, err := c.authService.Register(ctx.Context(), req)
	if err != nil {
		return response.SendResponse(ctx, responseStatusFromError(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusCreated, res)
}

func (c *AuthController) LoginUser(ctx fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, errx.ErrBadRequest.WithError(err))
	}
	res, err := c.authService.Login(ctx.Context(), req)
	if err != nil {
		return response.SendResponse(ctx, responseStatusFromError(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *AuthController) GetCurrentUser(ctx fiber.Ctx) error {
	user, _ := ctx.Locals("user").(dto.UserResponse)
	return response.SendResponse(ctx, fiber.StatusOK, dto.MeResponse{User: user})
}

func responseStatusFromError(err error) int {
	var reqErr *errx.RequestError
	if errors.As(err, &reqErr) {
		return reqErr.StatusCode
	}
	return fiber.StatusInternalServerError
}
