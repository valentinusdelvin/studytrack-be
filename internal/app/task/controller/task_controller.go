package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/valentinusdelvin/studytrack-be/domain/dto"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/internal/app/shared/httpx"
	"github.com/valentinusdelvin/studytrack-be/internal/app/task/service"
	"github.com/valentinusdelvin/studytrack-be/internal/middleware"
	"github.com/valentinusdelvin/studytrack-be/pkg/helpers/http/response"
)

type TaskController struct {
	taskService *service.TaskService
}

func NewTaskController(taskService *service.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

func (c *TaskController) GetTasks(ctx fiber.Ctx) error {
	res, err := c.taskService.GetTasks(ctx.Context(), middleware.UserID(ctx))
	if err != nil {
		return response.SendResponse(ctx, httpx.Status(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *TaskController) GetTaskByID(ctx fiber.Ctx) error {
	id, err := parseTaskID(ctx)
	if err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, err)
	}
	res, err := c.taskService.GetTaskByID(ctx.Context(), middleware.UserID(ctx), id)
	if err != nil {
		return response.SendResponse(ctx, httpx.Status(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusOK, res)
}
func (c *TaskController) CreateTask(ctx fiber.Ctx) error {
	var req dto.CreateTaskRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, errx.ErrBadRequest.WithError(err))
	}
	res, err := c.taskService.CreateTask(ctx.Context(), middleware.UserID(ctx), req)
	if err != nil {
		return response.SendResponse(ctx, httpx.Status(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusCreated, res)
}

func (c *TaskController) UpdateTask(ctx fiber.Ctx) error {
	id, err := parseTaskID(ctx)
	if err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, err)
	}
	var req dto.UpdateTaskRequest
	if err = ctx.Bind().Body(&req); err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, errx.ErrBadRequest.WithError(err))
	}
	res, err := c.taskService.UpdateTask(ctx.Context(), middleware.UserID(ctx), id, req)
	if err != nil {
		return response.SendResponse(ctx, httpx.Status(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func (c *TaskController) DeleteTask(ctx fiber.Ctx) error {
	id, err := parseTaskID(ctx)
	if err != nil {
		return response.SendResponse(ctx, fiber.StatusBadRequest, err)
	}
	res, err := c.taskService.DeleteTask(ctx.Context(), middleware.UserID(ctx), id)
	if err != nil {
		return response.SendResponse(ctx, httpx.Status(err), err)
	}
	return response.SendResponse(ctx, fiber.StatusOK, res)
}

func parseTaskID(ctx fiber.Ctx) (uint64, error) {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, errx.ErrBadRequest.WithMessage("ID tugas tidak valid.")
	}
	return id, nil
}
