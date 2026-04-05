package codegen

import (
	"os"
	"path/filepath"
)

// ScaffoldFiles creates empty source and test files for a rune in codeDir.
// It is idempotent — existing files are not overwritten.
func ScaffoldFiles(codeDir, shortName, ext string) error {
	if err := os.MkdirAll(codeDir, 0755); err != nil {
		return err
	}

	src := filepath.Join(codeDir, shortName+ext)
	test := filepath.Join(codeDir, testFilename(shortName, ext))

	for _, path := range []string{src, test} {
		if _, err := os.Stat(path); err == nil {
			continue // already exists
		}
		if err := os.WriteFile(path, nil, 0644); err != nil {
			return err
		}
	}
	return nil
}

// testFilename returns the conventional test filename for a given short name and extension.
func testFilename(shortName, ext string) string {
	switch ext {
	case ".go":
		return shortName + "_test.go"
	case ".py":
		return "test_" + shortName + ".py"
	default:
		return shortName + ".test" + ext
	}
}
