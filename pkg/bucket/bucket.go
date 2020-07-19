package bucket

import (
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v6"
)

// TODO wrap error

type Bucket struct {
	client              *minio.Client
	name                string
	presignedExpiration time.Duration
}

func NewBucket(client *minio.Client, name string, presignedExpiration time.Duration) (*Bucket, error) {
	if exist, errExist := client.BucketExists(name); errExist != nil || !exist {
		if err := client.MakeBucket(name, ""); err != nil {
			return nil, fmt.Errorf("client failed to make bucket: %w", err)
		}
	}

	return &Bucket{
		client:              client,
		name:                name,
		presignedExpiration: presignedExpiration,
	}, nil
}

func (b *Bucket) GetObject(objectName string) (*minio.Object, error) {
	object, err := b.client.GetObject(b.name, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("client failed to get object: %w", err)
	}

	return object, nil
}

func (b *Bucket) PutObject(objectName string, reader io.Reader, contentType string, objectSize int64) error {
	if _, err := b.client.PutObject(b.name, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	}); err != nil {
		return fmt.Errorf("client failed to put object: %w", err)
	}

	return nil
}

func (b *Bucket) FGetObject(objectName, path string) error {
	if err := b.client.FGetObject(b.name, objectName, path, minio.GetObjectOptions{}); err != nil {
		return fmt.Errorf("client failed to fget object")
	}

	return nil
}

func (b *Bucket) FPutObject(objectName, path string, contentType string) error {
	if _, err := b.client.FPutObject(b.name, objectName, path, minio.PutObjectOptions{
		ContentType: contentType,
	}); err != nil {
		return err
	}

	return nil
}

func (b *Bucket) PresignedGetObject(objectName string) (string, error) {
	presignedURL, err := b.client.PresignedGetObject(b.name, objectName, b.presignedExpiration, nil)
	if err != nil {
		return "", fmt.Errorf("client failed to presigned get object: %w", err)
	}

	return presignedURL.String(), nil
}

func DefaultPresignedExpiration() time.Duration {
	return time.Hour * 24
}

func (b *Bucket) StatObject(objectName string) (minio.ObjectInfo, error) {
	objectInfo, err := b.client.StatObject(b.name, objectName, minio.StatObjectOptions{})
	if err != nil {
		return minio.ObjectInfo{}, fmt.Errorf("client failed to stat object: %w", err)
	}

	return objectInfo, nil
}
