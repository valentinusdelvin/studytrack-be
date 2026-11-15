package main

import (
	"fmt"

	"github.com/valentinusdelvin/studytrack-be/internal/infra/database"
	"github.com/valentinusdelvin/studytrack-be/internal/infra/env"
	"github.com/valentinusdelvin/studytrack-be/internal/infra/server"
	"github.com/valentinusdelvin/studytrack-be/pkg/log"
)

func main() {
	db, err := database.NewPostgres(*env.AppEnv)
	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Fatal("failed to initialize database")
	}

	server.CommandHandler(db)

	httpServer := server.NewHTTPServer(db)
	httpServer.MountRoutes()
	httpServer.Start(fmt.Sprintf(":%s", env.AppEnv.AppPort))
}
