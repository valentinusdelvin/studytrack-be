package service

import (
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/valentinusdelvin/studytrack-be/domain/dto"
	"github.com/valentinusdelvin/studytrack-be/domain/entity"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/internal/app/auth/repository"
	"github.com/valentinusdelvin/studytrack-be/internal/app/shared/validation"
	appjwt "github.com/valentinusdelvin/studytrack-be/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	authRepo *repository.AuthRepository
	jwt      appjwt.CustomJwtInterface
}

func NewAuthService(authRepo *repository.AuthRepository, jwt appjwt.CustomJwtInterface) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		jwt:      jwt,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error) {
	if err := validation.Validate(req); err != nil {
		return dto.AuthResponse{}, err
	}
	if _, err := s.authRepo.FindUserByEmail(ctx, req.Email); err == nil {
		return dto.AuthResponse{}, errx.NewError(fiber.StatusConflict, "email_already_registered", "Email sudah terdaftar.")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.AuthResponse{}, errx.ErrInternalServer.WithError(err).WithLocation("auth.Register.findByEmail")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, errx.ErrInternalServer.WithError(err).WithLocation("auth.Register.hashPassword")
	}
	user := &entity.User{Name: strings.TrimSpace(req.Name), Email: strings.ToLower(strings.TrimSpace(req.Email)), PasswordHash: string(hash)}
	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		var requestErr *errx.RequestError
		if errors.As(err, &requestErr) {
			return dto.AuthResponse{}, err
		}
		return dto.AuthResponse{}, errx.ErrInternalServer.WithError(err).WithLocation("auth.Register.createUser")
	}
	return s.createAuthResponse(user)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	if err := validation.Validate(req); err != nil {
		return dto.AuthResponse{}, err
	}
	user, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.AuthResponse{}, errx.NewError(fiber.StatusUnauthorized, "invalid_credentials", "Email atau password salah.")
	}
	if err != nil {
		return dto.AuthResponse{}, errx.ErrInternalServer.WithError(err).WithLocation("auth.Login.findByEmail")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return dto.AuthResponse{}, errx.NewError(fiber.StatusUnauthorized, "invalid_credentials", "Email atau password salah.")
	}
	return s.createAuthResponse(user)
}

func (s *AuthService) createAuthResponse(user *entity.User) (dto.AuthResponse, error) {
	token, err := s.jwt.Create(user.ID, user.Name, user.Email)
	if err != nil {
		return dto.AuthResponse{}, errx.ErrInternalServer.WithError(err).WithLocation("auth.createToken")
	}
	return dto.AuthResponse{Token: token, User: dto.FormatToUserResponse(*user)}, nil
}
