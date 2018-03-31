package gcs

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	"github.com/twinj/uuid"
)

var (
	documents     = []string{".docx", ".docm", ".dot", ".dotx", ".dotm", ".html", ".txt", ".rtf", ".odt"}
	spreadsheets  = []string{".xlsx", ".xlsm", ".xlt", ".xltx", ".xltm", ".ods", ".csv", ".tsv", ".tab"}
	presentations = []string{".pptx", ".pptm", ".pps", ".ppsx", ".ppsm", ".pot", ".potx", ".potm"}
)

// TempName generates a random name
func TempName(prefix string, length int) string {
	randBytes := make([]byte, length)
	rand.Read(randBytes)
	return prefix + hex.EncodeToString(randBytes)
}

// GenUniqueFileKey returns a unique filename for safe file storage
func GenUniqueFileKey(filename string) string {
	u := uuid.NewV4()
	return fmt.Sprintf("%s%s", u, GetFileExt(filename))
}

// GetFileName returns the name of a file
func GetFileName(key string) string {
	return strings.TrimSuffix(key, GetFileExt(key))
}

// GetFileExt returns the extension of a file by filename(key)
func GetFileExt(key string) string {
	return filepath.Ext(key)
}

// GetMimeType returns mimetype of file
func GetMimeType(key string) string {
	return mime.TypeByExtension(GetFileExt(key))
}

// GetKind returns kind of file we are saving
func GetKind(key string) string {
	mimetype := GetMimeType(key)
	ext := GetFileExt(key)

	if strings.Contains(mimetype, "opendocument") {
		postfix := strings.Split(mimetype, "opendocument.")[1]

		if strings.Contains(postfix, "spreadsheet") {
			return fmt.Sprintf("spreadsheet/%s", ext)
		}

		if strings.Contains(postfix, "presentation") {
			return fmt.Sprintf("presentation/%s", ext)
		}

		return fmt.Sprintf("document/%s", ext)
	}

	if strings.Contains(mimetype, "officedocument") {
		postfix := strings.Split(mimetype, "officedocument.")[1]
		sep := strings.Split(postfix, ".")

		if strings.Contains(sep[0], "spreadsheet") {
			return fmt.Sprintf("spreadsheet/%s", ext)
		}

		if strings.Contains(sep[0], "presentation") {
			return fmt.Sprintf("presentation/%s", ext)
		}

		return fmt.Sprintf("document/%s", ext)
	}

	if mimetype == "application/msword" {
		return fmt.Sprintf("document/%s", ext)
	}

	if mimetype == "application/vnd.ms-excel" {
		return fmt.Sprintf("spreadsheet/%s", ext)
	}

	// if documents escape previous check
	if stringInSlice(ext, documents) {
		return fmt.Sprintf("document/%s", ext)
	}

	if stringInSlice(ext, spreadsheets) {
		return fmt.Sprintf("spreadsheet/%s", ext)
	}

	if stringInSlice(ext, presentations) {
		return fmt.Sprintf("presentation/%s", ext)
	}

	// make pdfs documents
	if mimetype == "application/pdf" {
		return fmt.Sprintf("document/%s", ext)
	}

	return mimetype
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
