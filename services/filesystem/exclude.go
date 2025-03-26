package filesystem

import "strings"

func IsExcludedFile(fileName string) bool {
	lowerFileName := strings.Trim(strings.ToLower(fileName), " ")
	if lowerFileName == ".git" || lowerFileName == ".svn" || lowerFileName == ".hg" ||
		lowerFileName == ".idea" || lowerFileName == ".vscode" || lowerFileName == "node_modules" ||
		lowerFileName == ".ds_store" {
		return true
	}
	return false
}
