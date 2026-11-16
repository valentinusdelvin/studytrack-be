package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/valentinusdelvin/studytrack-be/domain/entity"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/pkg/helpers/pgerror"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("lower(email) = ?", strings.ToLower(strings.TrimSpace(email))).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *entity.User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err == nil {
		return nil
	}

	if postgresErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		if mapped := pgerror.HandlePgError(*postgresErr, []pgerror.PgError{{
			Code:           pgerror.UniqueViolation,
			ConstraintName: "idx_users_email",
			Err: errx.NewError(
				fiber.StatusConflict,
				"email_already_registered",
				"Email sudah terdaftar.",
			),
		}}); mapped != nil {
			return mapped
		}
	}

	return err
}
