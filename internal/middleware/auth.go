package middleware

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/valentinusdelvin/studytrack-be/domain/dto"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/pkg/helpers/http/response"
	appjwt "github.com/valentinusdelvin/studytrack-be/pkg/jwt"
)

func RequireAuth(ctx fiber.Ctx) error {
	header := ctx.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return response.SendResponse(ctx, fiber.StatusUnauthorized, errx.ErrUnauthorized.WithMessage("Token tidak ditemukan. Silakan login terlebih dahulu."))
	}
	var claims appjwt.Claims
	if err := appjwt.Jwt.Decode(strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")), &claims); err != nil {
		return response.SendResponse(ctx, fiber.StatusUnauthorized, errx.ErrUnauthorized.WithMessage("Token tidak valid atau sudah kedaluwarsa.").WithError(err))
	}
	ctx.Locals("user", dto.UserResponse{ID: claims.UserID, Name: claims.Name, Email: claims.Email})
	ctx.Locals("user_id", claims.UserID)
	return ctx.Next()
}

func UserID(ctx fiber.Ctx) uint64 {
	if id, ok := ctx.Locals("user_id").(uint64); ok {
		return id
	}
	value, ok := ctx.Locals("user_id").(string)
	if !ok {
		return 0
	}
	id, _ := strconv.ParseUint(value, 10, 64)
	return id
}
