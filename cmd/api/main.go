package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Napat/mcpserver-demo/internal/router"
	"github.com/Napat/mcpserver-demo/pkg/database"
	"github.com/Napat/mcpserver-demo/pkg/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	// โหลดไฟล์ .env
	err := godotenv.Load("configs/temp/.env")
	if err != nil {
		log.Printf("Warning: .env file not found or invalid: %v", err)
	}

	// กำหนด Redis address ในกรณีที่ยังไม่มีตั้งค่าไว้
	if os.Getenv("REDIS_ADDR") == "" {
		// ตรวจสอบว่ารันใน Docker หรือไม่
		if os.Getenv("DOCKER_ENV") == "true" {
			// ถ้ารันใน Docker, ใช้ชื่อ service จาก docker-compose
			os.Setenv("REDIS_ADDR", "redis:6379")
		} else {
			// ถ้ารันใน local, ใช้ host.docker.internal ซึ่งชี้ไปที่ host machine
			os.Setenv("REDIS_ADDR", "host.docker.internal:6379")
		}
	}

	// สร้างเซิร์ฟเวอร์ Echo
	e := echo.New()

	// เพิ่ม middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ตั้งค่า CORS ให้ยอมรับการเชื่อมต่อจาก Frontend
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	allowedOrigins := []string{"*"} // ใช้ * เมื่อทำงานกับ nginx

	// ถ้ามีการกำหนดค่า CORS_ALLOWED_ORIGINS เฉพาะเจาะจง ให้ใช้ค่านั้นแทน
	if corsOrigins != "" && corsOrigins != "*" {
		allowedOrigins = strings.Split(corsOrigins, ",")
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}))

	// ลงทะเบียน validator
	validator.RegisterValidator(e)

	// เชื่อมต่อกับฐานข้อมูล
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// สร้าง logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// ตรวจสอบการเชื่อมต่อกับฐานข้อมูล
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database connection", zap.Error(err))
	}
	if err := sqlDB.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Info("Connected to database successfully")

	router.SetupRoutes(e, db, logger)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}

	// เริ่มเซิร์ฟเวอร์ API
	logger.Info("API Server is running", zap.String("url", "http://"+host+":"+port))
	e.Logger.Fatal(e.Start(host + ":" + port))
}
