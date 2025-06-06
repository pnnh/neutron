package filesystem

import (
	"os"
	"path/filepath"
	"strings"
)

func IsTextFile(extName string) bool {
	switch extName {
	case ".md", ".txt":
		return true
	}
	return false
}

func IsImageFile(fileName string) bool {
	extName := strings.Trim(strings.ToLower(filepath.Ext(fileName)), " ")
	switch extName {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return true
	}
	return false
}

func LowerExtName(fileName string) string {
	extName := strings.Trim(strings.ToLower(filepath.Ext(fileName)), " ")
	if extName == "" {
		return ""
	}
	return extName
}

func MkdirAll(path string) error {
	if path == "" {
		return nil
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	if err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return os.MkdirAll(filepath.Dir(path), os.ModePerm)
	}); err != nil {
		return err
	}
	return nil
}
