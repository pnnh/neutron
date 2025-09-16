package helpers

import (
	"path"
	"strings"
)

func GetMimeType(filePath string) string {
	extName := strings.ToLower(path.Ext(filePath))
	switch extName {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js", ".mjs", ".cjs":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".tiff", ".tif":
		return "image/tiff"
	case ".ico":
		return "image/vnd.microsoft.icon"
	case ".svg":
		return "image/svg+xml"
	case ".txt", ".md", ".markdown":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	case ".mp3":
		return "audio/mpeg"
	case ".mp4":
		return "video/mp4"
	case ".avi":
		return "video/x-msvideo"
	case ".mov":
		return "video/quicktime"
	case ".wmv":
		return "video/x-ms-wmv"
	case ".go", ".mod", ".sum":
		return "text/plain"
	case ".cs":
		return "text/plain"
	case ".java":
		return "text/plain"
	case ".py":
		return "text/plain"
	case ".rb":
		return "text/plain"
	case ".php":
		return "text/plain"
	case ".c", ".h":
		return "text/plain"
	case ".cpp", ".hpp", ".cc", ".cxx", ".hh", ".hxx", ".ixx":
		return "text/plain"
	case ".rs":
		return "text/plain"
	case ".ts", ".tsx", ".d.ts":
		return "text/plain"
	case ".yml", ".yaml":
		return "text/plain"
	case ".xml":
		return "application/xml"
	case ".toml":
		return "text/plain"
	case ".ini":
		return "text/plain"
	case ".env":
		return "text/plain"
	case ".dockerignore":
		return "text/plain"
	case ".gitignore":
		return "text/plain"
	case ".gitattributes":
		return "text/plain"
	case ".editorconfig":
		return "text/plain"
	case ".eslintignore", ".eslintrc", ".eslintrc.json", ".eslintrc.js":
		return "text/plain"
	case ".prettierrc", ".prettierrc.json", ".prettierrc.js":
		return "text/plain"
	case ".babelrc", ".babelrc.json", ".babelrc.js":
		return "text/plain"
	case ".stylelintrc", ".stylelintrc.json", ".stylelintrc.js":
		return "text/plain"
	case ".workspace":
		return "text/plain"
	case ".bazelrc", ".bazelignore", ".bzl", ".BUILD", ".bazel":
		return "text/plain"
	case ".gradle", ".gradlew", ".gradlew.bat":
		return "text/plain"
	case ".webp":
		return "image/webp"
	}
	fileName := strings.ToLower(path.Base(filePath))
	switch fileName {
	case "dockerfile":
		return "text/plain"
	case "makefile":
		return "text/plain"
	case "workspace", "build": // bazel
		return "text/plain"
	}
	return "application/octet-stream"
}

func IsTextFile(filePath string) bool {
	mimeType := GetMimeType(filePath)
	if strings.HasPrefix(mimeType, "text/") {
		return true
	}
	return false
}

func IsImageFile(filePath string) bool {
	mimeType := GetMimeType(filePath)
	if strings.HasPrefix(mimeType, "image/") {
		return true
	}
	return false
}
