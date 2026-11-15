package httpx

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
)

func Status(err error) int {
	var reqErr *errx.RequestError
	if errors.As(err, &reqErr) {
		return reqErr.StatusCode
	}
	return fiber.StatusInternalServerError
}
