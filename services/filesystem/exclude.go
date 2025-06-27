package filesystem

import (
	"path/filepath"
	"strings"
)

func IsExcludedFile(fileName string) bool {
	lowerFileName := strings.Trim(strings.ToLower(fileName), " ")
	if lowerFileName == ".git" || lowerFileName == ".svn" || lowerFileName == ".hg" ||
		lowerFileName == ".idea" || lowerFileName == ".vscode" || lowerFileName == "node_modules" ||
		lowerFileName == ".ds_store" {
		return true
	}
	return false
}

// IsHidden 实现判断指定路径是否是隐藏文件或文件夹的函数，兼容Linux、macOS
func IsHidden(path string) (bool, error) {
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") {
		return true, nil
	}
	return false, nil
}
