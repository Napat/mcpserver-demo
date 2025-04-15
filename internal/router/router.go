package router

import (
	"net/http"

	"github.com/Napat/mcpserver-demo/internal/handler"
	"github.com/Napat/mcpserver-demo/internal/repository"
	"github.com/Napat/mcpserver-demo/internal/service"
	"github.com/Napat/mcpserver-demo/pkg/cache"
	"github.com/Napat/mcpserver-demo/pkg/middleware"
	"github.com/Napat/mcpserver-demo/pkg/storage"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupRoutes ตั้งค่าเส้นทาง API
func SetupRoutes(e *echo.Echo, db *gorm.DB, logger *zap.Logger) {
	// สร้าง dependencies
	fileStorage, err := storage.NewMinioStorage()
	if err != nil {
		logger.Fatal("Failed to initialize MinIO storage", zap.Error(err))
	}

	// สร้าง Redis client
	redisClient, err := cache.NewRedisClient()
	if err != nil {
		logger.Fatal("Failed to initialize Redis client", zap.Error(err))
	}

	// สร้าง repositories ตาม Facade pattern (รวมการเข้าถึง database และ storage)
	userRepo := repository.NewUserRepository(db, fileStorage)
	noteRepo := repository.NewNoteRepository(db)
	visitorRepo := repository.NewVisitorRepository(redisClient)

	// สร้าง services
	userService := service.NewUserService(userRepo, logger)
	noteService := service.NewNoteService(noteRepo, logger)
	visitorService := service.NewVisitorService(visitorRepo, logger)

	// สร้าง handlers
	authHandler := handler.NewAuthHandler(userService, logger)
	userHandler := handler.NewUserHandler(userService, logger)
	noteHandler := handler.NewNoteHandler(noteService, logger)
	visitorHandler := handler.NewVisitorHandler(visitorService, logger)

	// API Routes
	api := e.Group("/api")

	// Health Check
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "OK",
			"message": "Server is running",
		})
	})

	// Public Routes
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	// Visitor Routes (Public)
	api.GET("/visitors", visitorHandler.GetVisitorCount)
	api.POST("/visitors", visitorHandler.IncrementVisitorCount)

	// Protected Routes
	user := api.Group("/me")
	user.Use(middleware.JWTMiddleware())
	user.GET("", userHandler.GetProfile)
	user.PUT("", userHandler.UpdateProfile)
	user.POST("/profile-image", userHandler.UpdateProfileImage)
	user.GET("/login-history", userHandler.GetLoginHistory)

	// Notes Routes (Protected)
	notes := api.Group("/notes")
	notes.Use(middleware.JWTMiddleware())
	notes.GET("", noteHandler.GetAllNotes)
	notes.GET("/:id", noteHandler.GetNote)
	notes.POST("", noteHandler.CreateNote)
	notes.PUT("/:id", noteHandler.UpdateNote)
	notes.DELETE("/:id", noteHandler.DeleteNote)

	// Admin Routes
	admin := api.Group("/admin")
	admin.Use(middleware.JWTMiddleware())
	admin.Use(middleware.AdminMiddleware)

	// TODO: Add admin routes for user management
}
