package server

import (
	"context"
	"flag"
	"os"

	"github.com/valentinusdelvin/studytrack-be/internal/infra/database"
	"gorm.io/gorm"
)

func CommandHandler(db *gorm.DB) {
	var migrationCmd string
	var seederCmd bool

	flag.StringVar(&migrationCmd, "m", "", "Migrate database 'up' or 'down'")
	flag.BoolVar(&seederCmd, "s", false, "Seed database")
	flag.Parse()

	if migrationCmd != "" {
		database.Migrate(db, migrationCmd)
		os.Exit(0)
	}

	if seederCmd {
		if err := database.Seed(context.Background(), db); err != nil {
			panic(err)
		}
		os.Exit(0)
	}
}
