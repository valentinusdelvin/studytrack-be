package database

import (
	"fmt"

	"github.com/valentinusdelvin/studytrack-be/internal/infra/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(env env.Env) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		env.DBHost,
		env.DBPort,
		env.DBUser,
		env.DBPass,
		env.DBName,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
