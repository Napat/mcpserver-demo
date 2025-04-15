package cache

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient คือตัวแทนสำหรับการเชื่อมต่อกับ Redis
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient คืนค่า instance ใหม่ของ Redis client
func NewRedisClient() (*RedisClient, error) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")
	db := 0

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// ทดสอบการเชื่อมต่อ
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		Client: client,
	}, nil
}

// Get ดึงข้อมูลจาก Redis
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Set เก็บข้อมูลใน Redis
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Incr เพิ่มค่า counter ใน Redis
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}
