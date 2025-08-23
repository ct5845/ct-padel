package main

import (
	"ct-go-web-starter/src/utils"
	"log/slog"
	"os"
)

func main() {
	if err := os.MkdirAll("build", 0755); err != nil {
		slog.Error("Failed to create build directory", "error", err)
		panic(err)
	}

	if err := copyAssets(); err != nil {
		panic(err)
	}

	if err := copyAlpineJS(); err != nil {
		panic(err)
	}

	if err := copyHTMX(); err != nil {
		panic(err)
	}
}

func copyAssets() error {
	// Copy assets from src/static to build/static
	err := utils.CopyDir("src/static", "build/static")
	if err != nil {
		slog.Error("Failed to copy assets", "error", err)
		return err
	}

	slog.Info("Copied assets to build/static", "src", "src/static", "dst", "build/static")
	return nil
}

func copyAlpineJS() error {
	// Copy Alpine.js from node_modules to build/static
	srcPath := "node_modules/alpinejs/dist/cdn.min.js"
	dstPath := "build/static/alpine.min.js"

	err := utils.CopyFile(srcPath, dstPath)
	if err != nil {
		slog.Error("Failed to copy Alpine.js", "error", err)
		return err
	}

	slog.Info("Copied Alpine.js to build/static", "src", srcPath, "dst", dstPath)
	return nil
}

func copyHTMX() error {
	// Copy HTMX from node_modules to build/static
	srcPath := "node_modules/htmx.org/dist/htmx.min.js"
	dstPath := "build/static/htmx.min.js"

	err := utils.CopyFile(srcPath, dstPath)
	if err != nil {
		slog.Error("Failed to copy HTMX", "error", err)
		return err
	}

	slog.Info("Copied HTMX to build/static", "src", srcPath, "dst", dstPath)
	return nil
}
