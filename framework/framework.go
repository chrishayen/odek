package framework

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed go/dispatch.go
var GoDispatch string

//go:embed ts/dispatch.ts
var TSDispatch string

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

// EnsureDispatchForLang writes the appropriate dispatcher for the given language.
func EnsureDispatchForLang(outputPath, language string) error {
	switch language {
	case "ts":
		dir := filepath.Join(outputPath, "dispatch")
		path := filepath.Join(dir, "dispatch.ts")
		if _, err := os.Stat(path); err == nil {
			return nil
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		return os.WriteFile(path, []byte(TSDispatch), 0644)
	default:
		return EnsureDispatch(outputPath)
	}
}
