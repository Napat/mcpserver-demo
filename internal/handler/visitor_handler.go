package handler

import (
	"net/http"

	"github.com/Napat/mcpserver-demo/internal/service"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// VisitorHandler จัดการ HTTP requests ที่เกี่ยวข้องกับผู้เข้าชม
type VisitorHandler struct {
	visitorService service.IVisitorService
	logger         *zap.Logger
}

// NewVisitorHandler สร้าง instance ใหม่ของ VisitorHandler
func NewVisitorHandler(visitorService service.IVisitorService, logger *zap.Logger) *VisitorHandler {
	return &VisitorHandler{
		visitorService: visitorService,
		logger:         logger,
	}
}

// IncrementVisitorCount เพิ่มและคืนค่าจำนวนผู้เข้าชม
func (h *VisitorHandler) IncrementVisitorCount(c echo.Context) error {
	ctx := c.Request().Context()

	count, err := h.visitorService.IncrementVisitorCount(ctx)
	if err != nil {
		h.logger.Error("Failed to increment visitor count", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to increment visitor count",
		})
	}

	return c.JSON(http.StatusOK, map[string]int64{
		"visitor_count": count,
	})
}

// GetVisitorCount คืนค่าจำนวนผู้เข้าชมปัจจุบัน
func (h *VisitorHandler) GetVisitorCount(c echo.Context) error {
	ctx := c.Request().Context()

	count, err := h.visitorService.GetVisitorCount(ctx)
	if err != nil {
		h.logger.Error("Failed to get visitor count", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get visitor count",
		})
	}

	return c.JSON(http.StatusOK, map[string]int64{
		"visitor_count": count,
	})
}
