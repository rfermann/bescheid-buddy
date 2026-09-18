package document_service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/cors"
)

type MinIOClient interface {
	BucketExists(ctx context.Context, bucket string) (bool, error)
	MakeBucket(ctx context.Context, bucket string, opts minio.MakeBucketOptions) error
	PresignedPutObject(ctx context.Context, bucket, object string, expires time.Duration) (*url.URL, error)
	SetBucketCors(ctx context.Context, bucketName string, corsConfig *cors.Config) error
	StatObject(ctx context.Context, bucket, object string, opts minio.StatObjectOptions) (minio.ObjectInfo, error)
	GetObject(ctx context.Context, bucket, object string, opts minio.GetObjectOptions) (*minio.Object, error)
	RemoveObject(ctx context.Context, bucket, object string, opts minio.RemoveObjectOptions) error
}

type Service struct {
	minioClient MinIOClient
	s3Config    config.S3Config
}

func New(minioClient MinIOClient, s3Config config.S3Config) *Service {
	return &Service{
		minioClient: minioClient,
		s3Config:    s3Config,
	}
}

func fileKey(caseID, documentID uuid.UUID) string {
	return fmt.Sprintf("%s/%s", caseID, documentID)
}

func (s *Service) PresignUpload(ctx context.Context, userID, caseID, documentID uuid.UUID) (string, error) {
	return s.PresignObjectUpload(ctx, userID, fileKey(caseID, documentID))
}

func (s *Service) PresignObjectUpload(ctx context.Context, userID uuid.UUID, objectKey string) (string, error) {
	bucket := userID.String()

	exists, err := s.minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return "", err
	}

	if !exists {
		err := s.minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return "", err
		}

		err = s.minioClient.SetBucketCors(ctx, bucket, &cors.Config{
			CORSRules: []cors.Rule{
				{
					AllowedOrigin: []string{s.s3Config.FrontendURL},
					AllowedMethod: []string{"PUT", "GET", "HEAD", "POST", "DELETE"},
					AllowedHeader: []string{"*"},
					MaxAgeSeconds: 3600,
				},
			},
		})
		if err != nil {
			return "", err
		}
	}

	presignedURL, err := s.minioClient.PresignedPutObject(ctx, bucket, objectKey, 15*time.Minute)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

var ErrObjectNotFound = errors.New("object not found in storage")

func (s *Service) StatDocumentSize(ctx context.Context, userID, caseID, documentID uuid.UUID) (int64, error) {
	return s.StatObjectSize(ctx, userID, fileKey(caseID, documentID))
}

func (s *Service) StatObjectSize(ctx context.Context, userID uuid.UUID, objectKey string) (int64, error) {
	bucket := userID.String()

	info, err := s.minioClient.StatObject(ctx, bucket, objectKey, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return 0, ErrObjectNotFound
		}
		return 0, err
	}

	return info.Size, nil
}

func (s *Service) DownloadDocument(ctx context.Context, userID, caseID, documentID uuid.UUID, maxBytes int64) ([]byte, error) {
	return s.DownloadObject(ctx, userID, fileKey(caseID, documentID), maxBytes)
}

func (s *Service) DownloadObject(ctx context.Context, userID uuid.UUID, objectKey string, maxBytes int64) ([]byte, error) {
	bucket := userID.String()

	obj, err := s.minioClient.GetObject(ctx, bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	data, err := io.ReadAll(io.LimitReader(obj, maxBytes+1))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) RemoveDocument(ctx context.Context, userID, caseID, documentID uuid.UUID) error {
	bucket := userID.String()
	return s.minioClient.RemoveObject(ctx, bucket, fileKey(caseID, documentID), minio.RemoveObjectOptions{})
}
