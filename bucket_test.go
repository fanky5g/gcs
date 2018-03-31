package gcs

import (
	"log"
	"os"
	"testing"
)

var (
	testBucketName = "silverbird-images"
	client         *GCloudStorageAgent
)

func TestMain(m *testing.M) {
	c, err := CreateClient("")

	if err != nil {
		log.Fatal(err)
	}

	client = c

	code := m.Run()
	os.Exit(code)
}

// func TestCreateBucket(t *testing.T) {
// 	bucket, err := client.CreateBucket(testBucketName)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, bucket)
//
// 	assert.IsType(t, &storage.BucketHandle{}, bucket)
//
// 	// upload file
// }

// func TestEmptyBucket(t *testing.T) {
// 	err := client.EmptyBucket(testBucketName)
// 	assert.NoError(t, err)
// }
//
// func TestDeleteBucket(t *testing.T) {
// 	err := client.DeleteBucket(testBucketName)
// 	assert.NoError(t, err)
// }
