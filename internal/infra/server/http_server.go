package server

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/valentinusdelvin/studytrack-be/domain/errx"
	"github.com/valentinusdelvin/studytrack-be/internal/app/auth/controller"
	authRepo "github.com/valentinusdelvin/studytrack-be/internal/app/auth/repository"
	authService "github.com/valentinusdelvin/studytrack-be/internal/app/auth/service"
	courseController "github.com/valentinusdelvin/studytrack-be/internal/app/course/controller"
	courseRepo "github.com/valentinusdelvin/studytrack-be/internal/app/course/repository"
	courseService "github.com/valentinusdelvin/studytrack-be/internal/app/course/service"
	taskController "github.com/valentinusdelvin/studytrack-be/internal/app/task/controller"
	taskRepo "github.com/valentinusdelvin/studytrack-be/internal/app/task/repository"
	taskService "github.com/valentinusdelvin/studytrack-be/internal/app/task/service"
	"github.com/valentinusdelvin/studytrack-be/internal/middleware"
	"github.com/valentinusdelvin/studytrack-be/pkg/helpers/http/response"
	appjwt "github.com/valentinusdelvin/studytrack-be/pkg/jwt"
	"github.com/valentinusdelvin/studytrack-be/pkg/log"
	"gorm.io/gorm"
)

type HTTPServer interface {
	Start(port string)
	MountRoutes()
	GetApp() *fiber.App
}

type httpServer struct {
	app *fiber.App
	db  *gorm.DB
}

func NewHTTPServer(db *gorm.DB) HTTPServer {
	config := fiber.Config{
		CaseSensitive: true,
		AppName:       "StudyTrack-BE",
		ServerHeader:  "StudyTrack-BE",
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			status := fiber.StatusInternalServerError
			var fiberErr *fiber.Error
			var requestErr *errx.RequestError
			switch {
			case errors.As(err, &requestErr):
				status = requestErr.StatusCode
			case errors.As(err, &fiberErr):
				status = fiberErr.Code
			}
			return response.SendResponse(ctx, status, err)
		},
	}

	app := fiber.New(config)
	app.Use(requestid.New())

	return &httpServer{
		app: app,
		db:  db,
	}
}

func (s *httpServer) GetApp() *fiber.App {
	return s.app
}

func (s *httpServer) Start(port string) {

	fiberListenConfig := fiber.ListenConfig{
		EnablePrefork: false,
	}

	err := s.app.Listen(port, fiberListenConfig)

	if err != nil {
		log.WithFields(log.Fields{"error": err.Error()}).Fatal("Error starting server")
	}
}

func (s *httpServer) MountRoutes() {
	s.app.Get("/api/health", func(ctx fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"status": "ok", "message": "StudyTrack backend is running."})
	})
	api := s.app.Group("/api")
	auth := controller.New(authService.New(authRepo.New(s.db), appjwt.Jwt))
	authRoutes := api.Group("/auth")
	authRoutes.Post("/register", auth.Register)
	authRoutes.Post("/login", auth.Login)
	authRoutes.Get("/me", middleware.RequireAuth, auth.Me)

	courses := courseController.New(courseService.New(courseRepo.New(s.db)))
	courseRoutes := api.Group("/courses", middleware.RequireAuth)
	courseRoutes.Get("/", courses.List)
	courseRoutes.Post("/", courses.Create)
	courseRoutes.Put("/:id", courses.Update)
	courseRoutes.Delete("/:id", courses.Delete)

	tasks := taskController.New(taskService.New(taskRepo.New(s.db)))
	taskRoutes := api.Group("/tasks", middleware.RequireAuth)
	taskRoutes.Get("/", tasks.List)
	taskRoutes.Get("/:id", tasks.Get)
	taskRoutes.Post("/", tasks.Create)
	taskRoutes.Put("/:id", tasks.Update)
	taskRoutes.Delete("/:id", tasks.Delete)
}
