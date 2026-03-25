package framework

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed go/dispatch.go
var GoDispatch string

// FS provides access to all framework files.
//
//go:embed go/*
var FS embed.FS

// EnsureDispatch writes dispatch.go into {outputPath}/dispatch/ if it doesn't already exist.
func EnsureDispatch(outputPath string) error {
	dir := filepath.Join(outputPath, "dispatch")
	path := filepath.Join(dir, "dispatch.go")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(GoDispatch), 0644)
}
