package gcs

import (
	"fmt"
	"mime"
	"path/filepath"

	"github.com/twinj/uuid"
)

// GenUniqueKey returns a unique filename for safe file storage
func GenUniqueKey(filename string) string {
	u := uuid.NewV4()
	return fmt.Sprintf("%s%s", u, GetExt(filename))
}

// GetExt returns the extension of a file by filename(key)
func GetExt(key string) string {
	return filepath.Ext(key)
}

// GetMimeType returns mimetype of file
func GetMimeType(key string) string {
	return mime.TypeByExtension(GetExt(key))
}
