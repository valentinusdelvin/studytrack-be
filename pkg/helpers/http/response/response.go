package response

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/internal/infra/env"
)

func SendResponse(
	ctx fiber.Ctx,
	code int,
	payload any,
) error {
	if code >= 400 {
		payload = formatErrorPayload(payload)
	}

	return ctx.Status(code).JSON(payload)
}

func formatErrorPayload(payload any) any {
	err, isError := payload.(error)
	if !isError {
		return payload
	}

	var reqErr *errx.RequestError
	if !errors.As(err, &reqErr) {
		return map[string]any{
			"error": err.Error(),
		}
	}

	if env.AppEnv.AppEnv != "production" {
		return map[string]any{
			"error": reqErr.Message,
		}
	}

	return map[string]any{"error": reqErr.Message}
}
