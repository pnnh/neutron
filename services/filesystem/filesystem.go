package filesystem

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

func ResolvePath(path string) (string, error) {
	resolvedPath := strings.ReplaceAll(path, "\\", "/")

	if strings.HasPrefix(resolvedPath, "file://") {
		resolvedPath = strings.Replace(resolvedPath, "file://", "", -1)
	}

	sysType := runtime.GOOS
	if strings.HasPrefix(resolvedPath, "work/") {
		dir, err := os.Getwd()
		if err != nil {
			return path, fmt.Errorf("ResolvePath Getwd3: %s", err)
		}
		resolvedPath = strings.Replace(resolvedPath, "work/", dir+"/", 1)
	} else if strings.HasPrefix(resolvedPath, "home/") {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return path, fmt.Errorf("ResolvePath UserHomeDir3: %s", err)
		}
		resolvedPath = strings.Replace(resolvedPath, "home/", userHomeDir+"/", 1)
	} else if strings.HasPrefix(resolvedPath, "root/") {
		if sysType == "windows" {
			diskDrivers := []string{"C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S",
				"T", "U", "V", "W", "X", "Y", "Z"}
			for _, diskDriver := range diskDrivers {
				if strings.HasPrefix(resolvedPath, "root/"+diskDriver+"/") {
					resolvedPath = strings.Replace(resolvedPath, "root/"+diskDriver, diskDriver+":", 1)
					break
				}
			}
		} else {
			resolvedPath = strings.Replace(resolvedPath, "root/", "/", 1)
		}
	}
	if sysType == "windows" {
		resolvedPath = strings.ReplaceAll(resolvedPath, "/", "\\")
	}

	return resolvedPath, nil
}

func CopyFile(src, dst string) (err error) {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func(sourceFile *os.File) {
		err = sourceFile.Close()
	}(sourceFile)

	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func(destinationFile *os.File) {
		err = destinationFile.Close()
	}(destinationFile)

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}
