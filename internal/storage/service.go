package storage

import (
	"colossus/pkg/bucket"
	"colossus/pkg/cache"
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/gabriel-vasile/mimetype"

	"github.com/rs/xid"
)

type Service struct {
	bucket *bucket.Bucket
	cache  *cache.Cache
}

func NewService(
	b *bucket.Bucket,
	c *cache.Cache,
) *Service {
	return &Service{
		bucket: b,
		cache:  c,
	}
}

func (s *Service) Upload(ctx context.Context, file multipart.File, size int64) (FileInfo, error) {
	guid := xid.New()
	id := guid.String()

	contentType, err := mimetype.DetectReader(file)
	if err != nil {
		return FileInfo{}, fmt.Errorf("mimetype failed to detect reader: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return FileInfo{}, fmt.Errorf("file failed to seek start: %w", err)
	}

	if err := s.bucket.PutObjectWithSize(id, file, contentType.String(), size); err != nil {
		return FileInfo{}, fmt.Errorf("bucket failed to put object: %w", err)
	}

	fileInfo := FileInfo{
		ID:          id,
		ContentType: contentType.String(),
		Extension:   contentType.Extension(),
	}

	if err := s.cache.SetJSON(ctx, id, fileInfo); err != nil {
		return FileInfo{}, fmt.Errorf("cache failed to set json: %w", err)
	}

	return fileInfo, nil
}

func (s *Service) Download(ctx context.Context, id string) (string, error) {
	url, err := s.bucket.PresignedGetObject(id)
	if err != nil {
		return "", fmt.Errorf("bucket failed to presigned get object: %w", err)
	}

	return url, nil
}
