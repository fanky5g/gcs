package gcs

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"path/filepath"

	"log"
	"sync"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
)

var (
	//ErrBodyEmpty returns when upload is called with an empty body
	ErrBodyEmpty = errors.New("Expected body to not be empty")
	//ErrKeyEmpty is an error that returns when an operation is supplied with an empty key
	ErrKeyEmpty = errors.New("Key cannot be empty")
)

// createFile creates a new file object on GCloud and returns
func (gcloudAgent *GCloudStorageAgent) createFile(r io.Reader, aclRules []storage.ACLRule, key, bucketName string) (*storage.ObjectAttrs, error) {
	bucket := gcloudAgent.Bucket(bucketName)
	ctx := context.Background()

	object := bucket.Object(key)
	wc := object.NewWriter(ctx)

	// modify aclRules if exists
	wc.ObjectAttrs.ACL = aclRules

	if _, err := io.Copy(wc, r); err != nil {
		return nil, err
	}

	if err := wc.Close(); err != nil {
		return nil, err
	}

	objectAttrs, err := object.Attrs(ctx)
	if err != nil {
		return nil, err
	}

	return objectAttrs, nil
}

// Upload takes a request object and an optional key parameter and returns an UploadOutput object
func (gcloudAgent *GCloudStorageAgent) Upload(body io.Reader, filename, bucketName string, aclRules []storage.ACLRule, opts *storage.SignedURLOptions) (*FileMetadata, error) {
	if body == nil {
		return nil, ErrBodyEmpty
	}

	reader, writer := io.Pipe()
	var size uint64

	go func() {
		defer writer.Close()
		s, err := io.Copy(writer, body)
		if err != nil {
			log.Fatal(err)
		}

		size = uint64(s)
	}()

	key := filepath.Join(filepath.Dir(filename), GenUniqueKey(filename))
	ret, err := gcloudAgent.createFile(reader, aclRules, key, bucketName)

	if err != nil {
		return nil, err
	}

	formatted, err := gcloudAgent.FormatFile(&File{
		Bucket:     ret.Bucket,
		Size:       size,
		ActualSize: uint64(ret.Size),
		Key:        key,
		Location:   ret.MediaLink,
		LastModified: ret.Updated,
	}, opts)

	if err != nil {
		return nil, err
	}

	return formatted, nil
}

// GetSignedURL returns downloadable file path for a private file
func (gcloudAgent *GCloudStorageAgent) GetSignedURL(bucketName, key string, opts *storage.SignedURLOptions) (*string, error) {
	if bucketName == "" {
		return nil, ErrBucketEmpty
	}

	if key == "" {
		return nil, ErrKeyEmpty
	}

	// opts := &storage.SignedURLOptions{
	// 	GoogleAccessID: accessID,
	// 	PrivateKey:     privKey,
	// 	Method:         "GET",
	// 	Expires:        time.Now().UTC().Add(time.Minute * 5),
	// }

	url, err := storage.SignedURL(bucketName, key, opts)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

// GetFile gets file and returns a stream
func (gcloudAgent *GCloudStorageAgent) GetFile(path string) (out io.ReadCloser, length int64, err error) {
	bucket, key, err := parsePath(path)
	if err != nil {
		return nil, 0, err
	}

	return gcloudAgent.getFile(bucket, key)
}

// getFile gets file object from GCloud and returns a readcloser
func (gcloudAgent *GCloudStorageAgent) getFile(bucketName, key string) (out io.ReadCloser, length int64, err error) {
	if bucketName == "" {
		return out, length, ErrBucketEmpty
	}

	if key == "" {
		return out, length, ErrKeyEmpty
	}

	bucket := gcloudAgent.Bucket(bucketName)
	ctx := context.Background()

	object := bucket.Object(key)
	o, err := object.NewReader(ctx)
	if err != nil {
		return out, length, err
	}

	return o, o.Size(), err
}

// DeleteMultiple deletes an array of s3 objects in a slice of *s3.ObjectIdentifier
func (gcloudAgent *GCloudStorageAgent) DeleteMultiple(bucketName string, keys []string) error {
	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Add(len(keys))
	for _, key := range keys {
		go func(b, k string) {
			errChan <- func() error {
				defer wg.Done()
				success, err := gcloudAgent.DeleteFile(b, k)
				if err != nil {
					return err
				}

				if !success {
					return fmt.Errorf("Failed to delete key %s", k)
				}
				return nil
			}()
		}(bucketName, key)
	}

	done := make(chan bool)
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case err := <-errChan:
		close(errChan)
		return err
	}
}

// DeleteFile deletes an s3 object and return a boolean of true if deleted
func (gcloudAgent *GCloudStorageAgent) DeleteFile(bucketName, key string) (bool, error) {
	if bucketName == "" {
		return false, ErrBucketEmpty
	}

	if key == "" {
		return false, ErrKeyEmpty
	}

	bucket := gcloudAgent.Bucket(bucketName)
	ctx := context.Background()

	object := bucket.Object(key)
	if err := object.Delete(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func parsePath(p string) (string, string, error) {
	u, err := url.Parse(p)
	if err != nil {
		return "", "", err
	}

	if !strings.HasPrefix(p, "gs://") {
		return "", "", fmt.Errorf("expected url to be in the form gs://<bucket>/key")
	}

	return u.Host, strings.TrimPrefix(strings.TrimSuffix(u.Path, "/"), "/"), nil
}

func targetPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Download downloads a file from gcs to specified location
func (gcloudAgent *GCloudStorageAgent) Download(url, dest string) error {
	bucket, key, err := parsePath(url)
	if err != nil {
		return err
	}

	r, size, err := gcloudAgent.getFile(bucket, key)
	if err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer f.Close()
	c, e := io.Copy(f, r)
	if e != nil {
		return e
	}

	if c != size {
		return fmt.Errorf("Failed: size mismatch. Expected download size of %v but got %v", size, c)
	}

	return nil
}

// ListObjects lists objects in bucket path sent
func (gcloudAgent *GCloudStorageAgent) ListObjects(bucket string, query *storage.Query) ([]FileMetadata, error) {
	bucketHandle := gcloudAgent.Bucket(bucket)
	it := bucketHandle.Objects(context.Background(), query)
	if it == nil {
		return []FileMetadata{}, nil
	}

	var objects []FileMetadata
	for {
			ret, err := it.Next()
			if err == iterator.Done {
					break
			}

			if err != nil {
				return objects, err
			}

			formatted, err := gcloudAgent.FormatFile(&File{
				Bucket:     ret.Bucket,
				Size:       uint64(ret.Size),
				ActualSize: uint64(ret.Size),
				Key:        ret.Name,
				Location:   ret.MediaLink,
				LastModified: ret.Updated,
			}, nil)

			if err != nil {
				return objects, err
			}

			objects = append(objects, *formatted)
	}

	return objects, nil
}