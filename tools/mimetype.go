package tools

import (
	"mime"
	"path/filepath"
)

func GetMimeType(filename string) string {
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	return mimeType
}
