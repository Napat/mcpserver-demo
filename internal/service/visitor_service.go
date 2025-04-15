package service

import (
	"context"

	"github.com/Napat/mcpserver-demo/internal/repository"
	"go.uber.org/zap"
)

// IVisitorService คือ interface สำหรับ Visitor Service
type IVisitorService interface {
	IncrementVisitorCount(ctx context.Context) (int64, error)
	GetVisitorCount(ctx context.Context) (int64, error)
}

// VisitorService ทำหน้าที่จัดการข้อมูลผู้เข้าชม
type VisitorService struct {
	visitorRepo repository.IVisitorRepository
	logger      *zap.Logger
}

// NewVisitorService สร้าง instance ใหม่ของ VisitorService
func NewVisitorService(visitorRepo repository.IVisitorRepository, logger *zap.Logger) *VisitorService {
	return &VisitorService{
		visitorRepo: visitorRepo,
		logger:      logger,
	}
}

// IncrementVisitorCount เพิ่มจำนวนผู้เข้าชม
func (s *VisitorService) IncrementVisitorCount(ctx context.Context) (int64, error) {
	count, err := s.visitorRepo.IncrementVisitorCount(ctx)
	if err != nil {
		s.logger.Error("Failed to increment visitor count", zap.Error(err))
		return 0, err
	}
	return count, nil
}

// GetVisitorCount อ่านจำนวนผู้เข้าชมปัจจุบัน
func (s *VisitorService) GetVisitorCount(ctx context.Context) (int64, error) {
	count, err := s.visitorRepo.GetVisitorCount(ctx)
	if err != nil {
		s.logger.Error("Failed to get visitor count", zap.Error(err))
		return 0, err
	}
	return count, nil
}
