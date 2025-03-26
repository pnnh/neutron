package filesystem

func IsTextFile(extName string) bool {
	switch extName {
	case ".md", ".txt":
		return true
	}
	return false
}

func IsImageFile(extName string) bool {
	switch extName {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return true
	}
	return false
}
