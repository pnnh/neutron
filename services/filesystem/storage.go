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
	fullTargetDir := filepath.Join(p.targetRootPath, string(os.PathSeparator), targetDir)
	err := os.MkdirAll(fullTargetDir, os.ModeDir)
	//path := `E:\Temp\5e31a3427bba91844defbe545ca0c0fa3109379a\main`
	//err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		return "", fmt.Errorf("MkdirAll: %w", err)
	}
	fullTargetPath := filepath.Join(p.targetRootPath, string(os.PathSeparator), targetPath)
	err = CopyFile(srcPath, fullTargetPath)
	if err != nil {
		return "", fmt.Errorf("CopyFile: %w", err)
	}

	return fullTargetPath, nil
}
