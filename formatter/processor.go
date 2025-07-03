package formatter

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ProcessDir(dir string, within bool, force bool, inspect bool, defaultDir string) {
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		realPath := path
		if inspect {
			newPath := ReorganizeFile(path, defaultDir)
			if newPath != "" {
				realPath = newPath
			}
		}
		ProcessFile(realPath, within, force, inspect)
		return nil
	})
}

func ReorganizeFile(path string, defaultDir string) string {
	dir := filepath.Dir(path)
	filename := strings.TrimSuffix(filepath.Base(path), ".md")
	if filepath.Base(path) == "index.md" {
		return path
	}
	targetDir := filepath.Join(dir, filename)
	os.MkdirAll(targetDir, os.ModePerm)
	newPath := filepath.Join(targetDir, "index.md")
	err := os.Rename(path, newPath)
	if err != nil {
		fmt.Printf("移动失败: %s\n", path)
		return path
	}
	if defaultDir != "" {
		entries, _ := os.ReadDir(defaultDir)
		for _, entry := range entries {
			src := filepath.Join(defaultDir, entry.Name())
			dest := filepath.Join(targetDir, entry.Name())
			data, err := os.ReadFile(src)
			if err == nil {
				os.WriteFile(dest, data, 0644)
			}
		}
	}
	fmt.Printf("已整理: %s -> %s\n", path, newPath)
	return newPath
}
