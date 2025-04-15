package handler

import (
	"net/http"

	"github.com/Napat/mcpserver-demo/internal/service"
	"github.com/Napat/mcpserver-demo/models"
	"github.com/Napat/mcpserver-demo/pkg/middleware"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoginRequest for login data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest for registration data
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Gender    string `json:"gender" validate:"required,oneof=male female other"`
}

// AuthHandler handles authentication
type AuthHandler struct {
	userService service.IUserService
	logger      *zap.Logger
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(userService service.IUserService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		logger:      logger,
	}
}

// Login จัดการการเข้าสู่ระบบ
func (h *AuthHandler) Login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to login", zap.Error(err))
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	// บันทึกประวัติการเข้าสู่ระบบ
	err = h.userService.RecordLogin(uint(user.ID), c.RealIP(), c.Request().UserAgent())
	if err != nil {
		h.logger.Error("Failed to record login history", zap.Error(err))
	}

	// สร้าง JWT token
	token, err := middleware.GenerateToken(uint(user.ID), user.Role)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// Register จัดการการลงทะเบียน
func (h *AuthHandler) Register(c echo.Context) error {
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// สร้างผู้ใช้ใหม่
	user := models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Gender:    req.Gender,
		Role:      models.RoleUser,
		Active:    true,
	}

	err := h.userService.Register(&user)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to register user")
	}

	// สร้าง JWT token
	token, err := middleware.GenerateToken(uint(user.ID), user.Role)
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	user.Password = "" // ไม่ส่งรหัสผ่านกลับไป
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
