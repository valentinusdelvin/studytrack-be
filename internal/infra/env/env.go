package env

import (
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/valentinusdelvin/studytrack-be/pkg/log"
)

type Env struct {
	AppEnv                  string        `mapstructure:"APP_ENV"`
	AppUrl                  string        `mapstructure:"APP_URL"`
	FEAppURL                string        `mapstructure:"FE_APP_URL"`
	AppPort                 string        `mapstructure:"APP_PORT"`
	DBHost                  string        `mapstructure:"DB_HOST"`
	DBPort                  string        `mapstructure:"DB_PORT"`
	DBUser                  string        `mapstructure:"DB_USER"`
	DBPass                  string        `mapstructure:"DB_PASS"`
	DBName                  string        `mapstructure:"DB_NAME"`
	DBMaxConnection         int           `mapstructure:"DB_MAX_CONNECTION"`
	DBMaxIdleConnection     int           `mapstructure:"DB_MAX_IDLE_CONNECTION"`
	DBMaxLifetimeConnection time.Duration `mapstructure:"DB_MAX_LIFETIME_CONNECTION"`
	DBMaxIdleTimeConnection time.Duration `mapstructure:"DB_MAX_IDLE_TIME_CONNECTION"`
	JwtSecretKey            string        `mapstructure:"JWT_SECRET_KEY"`
	JwtExpTime              time.Duration `mapstructure:"JWT_EXP_TIME"`
}

var AppEnv = getEnv()

func getEnv() *Env {
	env := &Env{}

	viper.SetConfigFile("./config/.env")

	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("[Env][GetEnv] failed to read config file, using env vars")
	}

	viper.AutomaticEnv()

	if err := viper.Unmarshal(env); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("[Env][GetEnv] failed to unmarshal to struct")
	}
	if env.JwtExpTime <= 0 {
		env.JwtExpTime = 7 * 24 * time.Hour
	}
	env.DBHost = strings.TrimSpace(env.DBHost)
	env.DBPort = strings.TrimSpace(env.DBPort)
	env.DBUser = strings.TrimSpace(env.DBUser)
	env.DBName = strings.TrimSpace(env.DBName)

	switch env.AppEnv {
	case "development":
		log.Info("Application is running on development mode")
	case "staging":
		log.Info("Application is running on staging mode")
	case "production":
		log.Info("Application is running on production mode")
	}

	return env
}

func IsProduction() bool {
	return getEnv().AppEnv == "production"
}
