package utils

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

func CopyFile(src, dst string) error {
	slog.Debug("Copying file", "src", src, "dst", dst)
	srcFile, err := os.Open(src)
	if err != nil {
		slog.Error("Failed to open source file", "src", src, "error", err)
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		slog.Error("Failed to create destination file", "dst", dst, "error", err)
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		slog.Error("Failed to copy file content", "error", err)
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	slog.Debug("File copied successfully", "src", src, "dst", dst)
	return nil
}
