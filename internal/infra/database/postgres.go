package database

import (
	"database/sql"
	"fmt"

	"github.com/valentinusdelvin/studytrack-be/internal/infra/env"
	"github.com/valentinusdelvin/studytrack-be/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormlogger "gorm.io/gorm/logger"
)

func NewPostgres(cfg env.Env) (*gorm.DB, *sql.DB) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBHost,
		cfg.DBPort,
		"disable",
	)

	gormLogLevel := gormlogger.Warn
	if env.IsProduction() {
		gormLogLevel = gormlogger.Error
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
	})

	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Fatal("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Fatal("failed to get sql.DB from gorm.DB")
	}

	// Connection Pool
	sqlDB.SetMaxOpenConns(env.AppEnv.DBMaxConnection)
	sqlDB.SetMaxIdleConns(env.AppEnv.DBMaxIdleConnection)
	sqlDB.SetConnMaxLifetime(env.AppEnv.DBMaxLifetimeConnection)
	sqlDB.SetConnMaxIdleTime(env.AppEnv.DBMaxIdleTimeConnection)

	return db, sqlDB
}
