package gcs

import (
	"strings"

	"github.com/rubenfonseca/fastimage"
	"cloud.google.com/go/storage"
)

// FormatFile formats a file and returns its metadata
func (gcloudAgent *GCloudStorageAgent) FormatFile(file *File, opts *storage.SignedURLOptions) (*FileMetadata, error) {
	mimetype := GetMimeType(file.Key)

	out := &FileMetadata{
		FileName: file.FileName,
		Key:      file.Key,
		MimeType: mimetype,
		Bucket:   file.Bucket,
		URL:      file.Location,
		FileSize: file.Size,
		LastModified: file.LastModified,
	}

	if strings.Contains(out.MimeType, "image") {
		// get signed url
		url, err := gcloudAgent.GetSignedURL(out.Bucket, out.Key, opts)
		if err != nil {
			return nil, err
		}

		imagetype, width, height, err := getImageMeta(*url)
		// Do nothing for failed requests
		if err == nil {
			out.Image = &Image{
				Width:     width,
				Height:    height,
				ImageType: imagetype,
			}
		}
	}

	return out, nil
}

func getImageMeta(url string) (string, uint32, uint32, error) {
	imagetype, size, err := fastimage.DetectImageType(url)
	if err != nil {
		return "", uint32(0), uint32(0), err
	}

	var itype string

	switch imagetype {
	case fastimage.JPEG:
		itype = "JPEG"
	case fastimage.PNG:
		itype = "PNG"
	case fastimage.GIF:
		itype = "GIF"
	}

	return itype, size.Width, size.Height, nil
}
