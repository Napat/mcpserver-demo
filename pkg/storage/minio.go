package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// IFileStorage interface สำหรับจัดการไฟล์
type IFileStorage interface {
	UploadFile(bucketName string, objectName string, file *multipart.FileHeader) (string, error)
	DeleteFile(bucketName string, objectName string) error
}

// MinioStorage struct สำหรับจัดการไฟล์ด้วย MinIO
type MinioStorage struct {
	client *minio.Client
}

// NewMinioStorage สร้าง instance ใหม่ของ MinioStorage
func NewMinioStorage() (*MinioStorage, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	useSSL := false

	// เชื่อมต่อกับ MinIO
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	return &MinioStorage{client: client}, nil
}

// UploadFile อัพโหลดไฟล์ไปยัง MinIO
func (s *MinioStorage) UploadFile(bucketName string, objectName string, file *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	// ตรวจสอบว่ามี bucket หรือไม่ ถ้าไม่มีให้สร้าง
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to create bucket: %w", err)
		}

		// ตั้งค่า policy เพื่อให้เข้าถึงได้แบบสาธารณะ
		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, bucketName)

		err = s.client.SetBucketPolicy(ctx, bucketName, policy)
		if err != nil {
			fmt.Printf("Warning: Failed to set bucket policy: %v\n", err)
			// ไม่ return error เพื่อให้ยังคงอัปโหลดไฟล์ได้
		}
	}

	// เปิดไฟล์
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// แก้ไข objectName ไม่ให้มีชื่อ bucket ซ้ำซ้อน
	// ลบคำว่า "profiles/" ออกจาก objectName ถ้ามี
	cleanObjectName := objectName
	if len(objectName) > len(bucketName)+1 && objectName[:len(bucketName)+1] == bucketName+"/" {
		cleanObjectName = objectName[len(bucketName)+1:]
	}

	// อัพโหลดไฟล์
	_, err = s.client.PutObject(ctx, bucketName, cleanObjectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// สร้าง URL สำหรับเข้าถึงไฟล์จากภายนอก
	// ให้ใช้ MINIO_PUBLIC_URL ถ้ามี มิฉะนั้นจะใช้ MINIO_ENDPOINT
	publicURL := os.Getenv("MINIO_PUBLIC_URL")
	if publicURL == "" {
		// ถ้าไม่มี MINIO_PUBLIC_URL ให้ใช้ MINIO_ENDPOINT แทน
		endpoint := os.Getenv("MINIO_ENDPOINT")
		if endpoint == "" {
			endpoint = "localhost:9000"
		}
		publicURL = fmt.Sprintf("http://%s", endpoint)
	}

	fileURL := fmt.Sprintf("%s/%s/%s", publicURL, bucketName, cleanObjectName)
	return fileURL, nil
}

// DeleteFile ลบไฟล์จาก MinIO
func (s *MinioStorage) DeleteFile(bucketName string, objectName string) error {
	ctx := context.Background()

	err := s.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetDefaultProfileImage คืนค่า URL ของรูปโปรไฟล์ดีฟอลต์
func (s *MinioStorage) GetDefaultProfileImage() string {
	// ให้ใช้ MINIO_PUBLIC_URL ถ้ามี มิฉะนั้นจะใช้ MINIO_ENDPOINT
	publicURL := os.Getenv("MINIO_PUBLIC_URL")
	if publicURL == "" {
		// ถ้าไม่มี MINIO_PUBLIC_URL ให้ใช้ MINIO_ENDPOINT แทน
		endpoint := os.Getenv("MINIO_ENDPOINT")
		if endpoint == "" {
			endpoint = "localhost:9000"
		}
		publicURL = fmt.Sprintf("http://%s", endpoint)
	}

	return fmt.Sprintf("%s/profiles/default.png", publicURL)
}
