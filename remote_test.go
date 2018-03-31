package gcs

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	// "github.com/rubenfonseca/fastimage"
)

var (
	tempFileKey string
	testFile1   = "gs://silverbird-images/hero_alicia_mobile_small.png"
	testFile2   = "gs://silverbird-images/myAvatar.png"
)

func TestDownload(t *testing.T) {
	dest, err := filepath.Abs("./alicia.png")
	if assert.NoError(t, err) {
		defer os.Remove(dest)
		err = client.Download(testFile1, dest)
		assert.NoError(t, err)
	}
}

// func upfile(isChunked bool) ([]*types.FileMetadata, *string, error) {
// 	tempFileKey = "alicia.mp4"
// 	fp, err := filepath.Abs(fmt.Sprintf("../data/static/testdata/%s", tempFileKey))
//
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	file, err := os.Open(fp)
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	defer file.Close()
//
// 	out, _, err := client.Upload(file, tempFileKey, "silverbird-raw-files", isChunked)
// 	if err != nil {
// 		return nil, nil, err
// 	}
//
// 	return out, nil, err
// }

// func TestUploadFileChunked(t *testing.T) {
// 	out, _, err := upfile(true)
//
// 	if err != nil {
// 		t.Errorf("Upload failed %s", err)
// 	}
//
// 	assert.IsType(t, []*types.FileMetadata{}, out)
// 	if out.Image.ImageType != "PNG" {
// 		t.Errorf("Expected imagetype to be %s but got %s", fastimage.JPEG, out.Image.ImageType)
// 	}
//
// 	if out.Image.Width == 0 {
// 		t.Error("Expected width to not be 0")
// 	}
// }

// func TestUploadFile(t *testing.T) {
// 	out, _, err := upfile(false)
//
// 	if err != nil {
// 		t.Errorf("Upload failed %s", err)
// 	}
//
// 	assert.IsType(t, []*types.FileMetadata{}, out)
// 	uploaded := out[0]
// 	t.Log(uploaded.URL)
// 	t.Log(uploaded.Key)
// 	t.Log(uploaded.Bucket)
// }

// func TestGetFile(t *testing.T) {
// 	r, length, err := Agent.GetFile(tempBucket, tempFileKey)
//
// 	if err != nil {
// 		t.Errorf("GetFile returned an error %s", err)
// 	}
//
// 	defer r.Close()
//
// 	var b bytes.Buffer
// 	writer := bufio.NewWriter(&b)
// 	bytesWritten, err := io.Copy(writer, r)
// 	// writer.Flush()
//
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	if bytesWritten != length {
// 		t.Errorf("Expected content length to equal %d but got %d", length, bytesWritten)
// 	}
// }

// func TestGetSignedURL(t *testing.T) {
//
// }

// func TestDeleteFile(t *testing.T) {
// 	deleted, err := Agent.DeleteFile(tempBucket, tempFileKey)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	if deleted != true {
// 		t.Errorf("Expected deleted to equal true but got %d", deleted)
// 	}
// }

// func TestEmptyS3Bucket(t *testing.T) {
// 	// reupload file
// 	_, err := upfile()
//
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	_, err = Agent.EmptyBucket("")
//
// 	if err != ErrBucketEmpty {
// 		t.Error(err)
// 	}
//
// 	out, err := Agent.EmptyBucket(tempBucket)
//
// 	if err != nil {
// 		t.Errorf("Empty bucket operation failed with error %s", err)
// 	}
//
// 	if len(out.Errors) != 0 {
// 		for _, e := range out.Errors {
// 			t.Error(e)
// 		}
// 	}
//
// 	if len(out.Deleted) != 1 {
// 		t.Errorf("Expected deleted object length to equal 1 but got %d", len(out.Deleted))
// 	}
// }

// func TestDeleteS3Bucket(t *testing.T) {
// 	_, err := Agent.DeleteS3Bucket(tempBucket)
//
// 	if err != nil {
// 		t.Error("Bucket delete failed %s", err)
// 	}
// }
