package controller

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/valentinusdelvin/studytrack-be/domain/dto"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/internal/app/auth/service"
	"github.com/valentinusdelvin/studytrack-be/pkg/helpers/http/response"
)

type Controller struct{ service *service.Service }

func New(service *service.Service) *Controller { return &Controller{service: service} }

func (c *Controller) Register(ctx fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, errx.ErrBadRequest.WithError(err))
	}
	res, err := c.service.Register(ctx.Context(), req)
	if err != nil {
		return response.SendResponse(ctx, errorStatus(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusCreated, res)
}

func (c *Controller) Login(ctx fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, errx.ErrBadRequest.WithError(err))
	}
	res, err := c.service.Login(ctx.Context(), req)
	if err != nil {
		return response.SendResponse(ctx, errorStatus(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *Controller) Me(ctx fiber.Ctx) error {
	user, _ := ctx.Locals("user").(dto.UserResponse)
	return response.SendResponse(ctx, fiber.StatusOK, dto.MeResponse{User: user})
}

func errorStatus(err error) int {
	var reqErr *errx.RequestError
	if errors.As(err, &reqErr) {
		return reqErr.StatusCode
	}
	return fiber.StatusInternalServerError
}
