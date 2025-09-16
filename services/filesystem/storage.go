package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type FilePorter struct {
	targetRootPath string
}

func NewFilePorter(targetPath string) (*FilePorter, error) {

	resolvedPath, err := ResolvePath(targetPath)
	if err != nil {
		logrus.Fatalln("NewFilePorter解析路径失败", err)
		return nil, fmt.Errorf("NewFilePorter解析路径失败")
	}
	return &FilePorter{targetRootPath: resolvedPath}, nil
}

func (p *FilePorter) CopyFile(srcPath, targetPath string) (string, error) {
	targetDir := filepath.Dir(targetPath)
	fullTargetDir := p.targetRootPath + "/" + targetDir
	err := os.MkdirAll(fullTargetDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("MkdirAll: %w", err)
	}
	fullTargetPath := p.targetRootPath + "/" + targetPath
	err = CopyFile(srcPath, fullTargetPath)
	if err != nil {
		return "", fmt.Errorf("CopyFile: %w", err)
	}

	return fullTargetPath, nil
}
