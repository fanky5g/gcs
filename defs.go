package gcs

import (
	"encoding/json"
	"errors"
	"io"

	"database/sql/driver"

	"cloud.google.com/go/storage"
)

// Uploader holds characteristics every upload object must implement
type Uploader interface {
	Upload(body io.Reader, data interface{}) (*FileMetadata, error)
}

// GCloudStorageAgent contains useful abstractions for file storage and manipulations on gcloud
type GCloudStorageAgent struct {
	*storage.Client
	GoogleAccessID string
	PrivateKey     string
	ProjectID      string
}

// File encloses file properties
type File struct {
	Bucket     string
	Key        string
	FileName   string
	Location   string
	Size       uint64
	ActualSize uint64
	AuthorID   uint8
}

//FileMetadata type returns File metadata for uploaded files
type FileMetadata struct {
	FileName string `json:"fileName"`
	Key      string `json:"key"`
	FileSize uint64 `json:"size"`
	MimeType string `json:"mime"`
	Bucket   string `json:"bucket"`
	URL      string `json:"url"`

	// image fields, omitted if not set
	Image *Image `json:"image" sql:"type:jsonb"`
}

// Image holds fields for image dimensions and type
type Image struct {
	Width     uint32 `json:"width, omitempty"  validate:"false"  editable:"true"`
	Height    uint32 `json:"height, omitempty"  validate:"false"  editable:"true"`
	ImageType string `json:"imageType, omitempty"  validate:"false"  editable:"true"`
}

// Value serializes image props from database into JSON
func (image Image) Value() (driver.Value, error) {
	j, err := json.Marshal(image)
	return j, err
}

// Scan serializes JSON into bytes for storage
func (image *Image) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed")
	}

	err := json.Unmarshal(source, &image)
	if err != nil {
		return err
	}

	return nil
}
