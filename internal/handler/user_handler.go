//go:generate mockgen -source=./user_handler.go -destination=./mocks/mock_user_handler.go -package=mocks

package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Napat/mcpserver-demo/internal/service"
	"github.com/Napat/mcpserver-demo/pkg/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ProfileUpdateRequest สำหรับข้อมูลการอัพเดทโปรไฟล์
type ProfileUpdateRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Gender    string `json:"gender" validate:"required,oneof=male female other"`
}

// UserHandler จัดการเกี่ยวกับ user endpoints
type UserHandler struct {
	userService service.IUserService
	logger      *zap.Logger
}

// NewUserHandler สร้าง instance ใหม่ของ UserHandler
func NewUserHandler(userService service.IUserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetProfile ดึงข้อมูลโปรไฟล์ของผู้ใช้ปัจจุบัน
func (h *UserHandler) GetProfile(c echo.Context) error {
	// Debug token claims
	claims, ok := c.Get("user").(jwt.MapClaims)
	if ok {
		fmt.Println("Token claims:", claims)
	} else {
		fmt.Println("No token claims found or wrong type")
	}

	// ดึงข้อมูลผู้ใช้จาก context
	userID := middleware.GetUserIDFromToken(c)
	fmt.Println("User ID from token:", userID)
	if userID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// ดึงข้อมูลผู้ใช้จากฐานข้อมูล
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		h.logger.Error("Failed to get user profile", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user profile")
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateProfile อัพเดทข้อมูลโปรไฟล์ของผู้ใช้
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	// ดึงข้อมูลผู้ใช้จาก context
	userID := middleware.GetUserIDFromToken(c)
	if userID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// ดึงข้อมูลจาก request
	req := new(ProfileUpdateRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// ดึงข้อมูลผู้ใช้ปัจจุบัน
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		h.logger.Error("Failed to get user for update", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update profile")
	}

	// อัพเดทข้อมูล
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Gender = req.Gender // Assign string directly

	if err := h.userService.UpdateProfile(user); err != nil {
		h.logger.Error("Failed to update profile", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update profile")
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateProfileImage อัพเดทรูปโปรไฟล์
func (h *UserHandler) UpdateProfileImage(c echo.Context) error {
	// ดึงข้อมูลผู้ใช้จาก context
	userID := middleware.GetUserIDFromToken(c)
	if userID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// ดึงไฟล์จาก form
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid image file")
	}

	// อัพโหลดไฟล์
	imageURL, err := h.userService.UpdateProfileImage(userID, file)
	if err != nil {
		h.logger.Error("Failed to update profile image", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update profile image")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"image_url": imageURL,
	})
}

// GetLoginHistory ดึงประวัติการเข้าสู่ระบบ
func (h *UserHandler) GetLoginHistory(c echo.Context) error {
	// ดึงข้อมูลผู้ใช้จาก context
	userID := middleware.GetUserIDFromToken(c)
	if userID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// ดึงค่า limit จาก query parameter
	limitStr := c.QueryParam("limit")
	limit := 10 // ค่าดีฟอลต์
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// ดึงประวัติการเข้าสู่ระบบ
	history, err := h.userService.GetLoginHistory(userID, limit)
	if err != nil {
		h.logger.Error("Failed to get login history", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get login history")
	}

	return c.JSON(http.StatusOK, history)
}
