package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/valentinusdelvin/studytrack-be/domain/dto"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/internal/app/course/service"
	"github.com/valentinusdelvin/studytrack-be/internal/app/shared/httpx"
	"github.com/valentinusdelvin/studytrack-be/internal/middleware"
	"github.com/valentinusdelvin/studytrack-be/pkg/helpers/http/response"
	"strconv"
)

type Controller struct{ service *service.Service }

func New(service *service.Service) *Controller { return &Controller{service: service} }
func (c *Controller) List(ctx fiber.Ctx) error {
	res, err := c.service.List(ctx.Context(), middleware.UserID(ctx))
	if err != nil {
		return response.SendResponse(ctx, httpx.Status(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusOK, res)
}
func (c *Controller) Create(ctx fiber.Ctx) error {
	var req dto.CreateCourseRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, errx.ErrBadRequest.WithError(err))
	}
	res, err := c.service.Create(ctx.Context(), middleware.UserID(ctx), req)
	if err != nil {
		return response.SendResponse(ctx, httpx.Status(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusCreated, res)
}
func (c *Controller) Update(ctx fiber.Ctx) error {
	id, err := pathID(ctx)
	if err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, err)
	}
	var req dto.UpdateCourseRequest
	if err = ctx.Bind().Body(&req); err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, errx.ErrBadRequest.WithError(err))
	}
	res, err := c.service.Update(ctx.Context(), middleware.UserID(ctx), id, req)
	if err != nil {
		return response.SendResponse(ctx, httpx.Status(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusOK, res)
}
func (c *Controller) Delete(ctx fiber.Ctx) error {
	id, err := pathID(ctx)
	if err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, err)
	}
	res, err := c.service.Delete(ctx.Context(), middleware.UserID(ctx), id)
	if err != nil {
		return response.SendResponse(ctx, httpx.Status(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusOK, res)
}
func pathID(ctx fiber.Ctx) (uint64, error) {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, errx.ErrBadRequest.WithMessage("ID mata kuliah tidak valid.")
	}
	return id, nil
}
