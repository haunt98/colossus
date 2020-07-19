package storage

import (
	"colossus/pkg/bucket"
	"colossus/pkg/cache"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"

	"go.uber.org/zap"

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

func (s *Service) Upload(ctx context.Context, sugar *zap.SugaredLogger, fileHeader *multipart.FileHeader) (FileInfo, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return FileInfo{}, fmt.Errorf("file header failed to open: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			sugar.Error(err)
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return FileInfo{}, fmt.Errorf("io failed to copy: %w", err)
	}
	checksum := hex.EncodeToString(h.Sum(nil))

	if cachedID, err := s.cache.GetString(ctx, checksumPrefix(checksum)); err == nil {
		var cachedFileInfo FileInfo
		if cachedErr := s.cache.GetJSON(ctx, cachedID, &cachedFileInfo); cachedErr != nil {
			sugar.Warn("checksum exist but id not exist", "checksum", checksum, "id", cachedID)

			return s.uploadNew(ctx, fileHeader, file, checksum)
		}

		return cachedFileInfo, nil
	}

	return s.uploadNew(ctx, fileHeader, file, checksum)
}

func (s *Service) uploadNew(ctx context.Context, fileHeader *multipart.FileHeader, file multipart.File,
	checksum string) (FileInfo, error) {
	guid := xid.New()
	id := guid.String()

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return FileInfo{}, fmt.Errorf("file failed to seek start: %w", err)
	}

	contentType, err := mimetype.DetectReader(file)
	if err != nil {
		return FileInfo{}, fmt.Errorf("mimetype failed to detect reader: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return FileInfo{}, fmt.Errorf("file failed to seek start: %w", err)
	}

	if err := s.bucket.PutObject(id, file, contentType.String(), fileHeader.Size); err != nil {
		return FileInfo{}, fmt.Errorf("bucket failed to put object: %w", err)
	}

	var fileInfo = FileInfo{
		ID:          id,
		ContentType: contentType.String(),
		Extension:   contentType.Extension(),
		Size:        fileHeader.Size,
		Checksum:    checksum,
	}

	if err := s.cache.SetJSON(ctx, id, fileInfo); err != nil {
		return FileInfo{}, fmt.Errorf("cache failed to set json: %w", err)
	}

	if err := s.cache.SetString(ctx, checksumPrefix(checksum), id); err != nil {
		return FileInfo{}, fmt.Errorf("cache failed to set string: %w", err)
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

func checksumPrefix(checksum string) string {
	return "checksum:" + checksum
}
