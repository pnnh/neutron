package filesystem

import (
	"path/filepath"
	"strings"
)

func IsImageFile(path string) bool {
	extName := strings.ToLower(filepath.Ext(path))
	switch extName {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return true
	}
	return false
}
