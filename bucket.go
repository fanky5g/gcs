package gcs

import (
	"errors"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

var (
	//ErrBucketEmpty is an error that returns when an operation is supplied with an empty bucket name
	ErrBucketEmpty = errors.New("BucketName must be a non-empty string")
)

// CreateBucket creates a new GCloud Bucket
func (gcloudAgent *GCloudStorageAgent) CreateBucket(BucketName string) (*storage.BucketHandle, error) {
	if BucketName == "" {
		return nil, ErrBucketEmpty
	}

	bucket := gcloudAgent.Bucket(BucketName)
	ctx := context.Background()

	if err := bucket.Create(ctx, gcloudAgent.ProjectID, nil); err != nil {
		return nil, err
	}

	return bucket, nil
}

// DeleteBucket removes an existing GCloud bucket
func (gcloudAgent *GCloudStorageAgent) DeleteBucket(BucketName string) error {
	if BucketName == "" {
		return ErrBucketEmpty
	}

	bucket := gcloudAgent.Bucket(BucketName)
	ctx := context.Background()

	if err := bucket.Delete(ctx); err != nil {
		return err
	}

	return nil
}

// EmptyBucket empties the specified gcloud bucket
func (gcloudAgent *GCloudStorageAgent) EmptyBucket(BucketName string) error {
	if BucketName == "" {
		return ErrBucketEmpty
	}

	bucket := gcloudAgent.Bucket(BucketName)
	ctx := context.Background()

	objects := bucket.Objects(ctx, nil)

	for {
		object, err := objects.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return err
		}

		handle := bucket.Object(object.Name)
		if err := handle.Delete(ctx); err != nil {
			return err
		}
	}

	return nil
}
