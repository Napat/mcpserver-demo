package repository

import (
	"context"
	"strconv"

	"github.com/Napat/mcpserver-demo/pkg/cache"
	"github.com/go-redis/redis/v8"
)

const (
	// visitorCountKey คือ key ที่ใช้เก็บจำนวนผู้เข้าชมใน Redis
	visitorCountKey = "visitor:count"
)

// IVisitorRepository คือ interface สำหรับ Visitor Repository
type IVisitorRepository interface {
	IncrementVisitorCount(ctx context.Context) (int64, error)
	GetVisitorCount(ctx context.Context) (int64, error)
}

// VisitorRepository ทำหน้าที่จัดการข้อมูลผู้เข้าชม
type VisitorRepository struct {
	redisClient *cache.RedisClient
}

// NewVisitorRepository สร้าง instance ใหม่ของ VisitorRepository
func NewVisitorRepository(redisClient *cache.RedisClient) *VisitorRepository {
	return &VisitorRepository{
		redisClient: redisClient,
	}
}

// IncrementVisitorCount เพิ่มจำนวนผู้เข้าชมแล้วคืนค่าจำนวนปัจจุบัน
func (r *VisitorRepository) IncrementVisitorCount(ctx context.Context) (int64, error) {
	// เพิ่มค่าจำนวนผู้เข้าชม
	count, err := r.redisClient.Incr(ctx, visitorCountKey)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetVisitorCount คืนค่าจำนวนผู้เข้าชมปัจจุบัน
func (r *VisitorRepository) GetVisitorCount(ctx context.Context) (int64, error) {
	// อ่านค่าจำนวนผู้เข้าชม
	val, err := r.redisClient.Get(ctx, visitorCountKey)
	if err == redis.Nil {
		// ถ้าไม่มีค่า เริ่มต้นที่ 0 และบันทึกลง Redis
		err = r.redisClient.Set(ctx, visitorCountKey, 0, 0)
		if err != nil {
			return 0, err
		}
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	// แปลงค่าจาก string เป็น int64
	count, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}
