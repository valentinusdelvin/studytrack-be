package database

import (
	"fmt"

	"github.com/valentinusdelvin/studytrack-be/domain/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, command string) {
	migrator := db.Migrator()
	tables := []any{
		&entity.User{},
		&entity.Course{},
		&entity.Task{},
	}

	var err error
	if command == "up" {
		err = migrator.AutoMigrate(tables...)
	}

	if command == "down" {
		err = migrator.DropTable(tables...)
	}

	if err != nil {
		panic(err)
	}

	fmt.Printf("Migration %s completed successfully\n", command)
}
