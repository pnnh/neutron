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
	case ".awebp":
		return "image/webp"
	case ".heic", ".heif":
		return "image/heic"
	case ".avif":
		return "image/avif"
	case ".jsonc":
		return "application/json"
	case ".tsv":
		return "text/tab-separated-values"
	case ".csv":
		return "text/csv"
	case ".ics":
		return "text/calendar"
	case ".svgz":
		return "image/svg+xml"
	case ".log":
		return "text/plain"
	case ".rst":
		return "text/plain"
	case ".tex":
		return "text/plain"
	case ".lhs":
		return "text/plain"
	case ".clj", ".cljs", ".cljc", ".edn":
		return "text/plain"
	case ".fs", ".fsi", ".fsx":
		return "text/plain"
	case ".vb":
		return "text/plain"
	case ".vbs":
		return "text/plain"
	case ".ps1":
		return "text/plain"
	case ".psm1":
		return "text/plain"
	case ".sh":
		return "text/plain"
	case ".bash":
		return "text/plain"
	case ".zsh":
		return "text/plain"
	case ".fish":
		return "text/plain"
	case ".ksh":
		return "text/plain"
	case ".csh":
		return "text/plain"
	case ".tcsh":
		return "text/plain"
	case ".lua":
		return "text/plain"
	case ".r":
		return "text/plain"
	case ".jl":
		return "text/plain"
	case ".groovy":
		return "text/plain"
	case ".makefile":
		return "text/plain"
	case ".mk":
		return "text/plain"
	case ".cmake":
		return "text/plain"
	case ".dockerfile":
		return "text/plain"
	case ".conf":
		return "text/plain"
	case ".config":
		return "text/plain"
	case ".props":
		return "text/plain"
	case ".properties":
		return "text/plain"
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
	switch mimeType {
	case "application/json", "application/xml", "application/javascript":
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
