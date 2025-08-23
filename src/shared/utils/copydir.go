package utils

import (
	"fmt"
	"log/slog"
	"os"
)

func CopyDir(src, dst string) error {
	slog.Debug("Copying directory", "src", src, "dst", dst)
	err := os.MkdirAll(dst, 0755)
	if err != nil {
		slog.Error("Failed to create destination directory", "dst", dst, "error", err)
		return fmt.Errorf("failed to create destination directory %s: %w", dst, err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		slog.Error("Failed to read source directory", "src", src, "error", err)
		return fmt.Errorf("failed to read source directory %s: %w", src, err)
	}

	for _, entry := range entries {
		srcPath := fmt.Sprintf("%s/%s", src, entry.Name())
		dstPath := fmt.Sprintf("%s/%s", dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	slog.Debug("Directory copied successfully", "src", src, "dst", dst)
	return nil
}
