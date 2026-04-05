package codegen

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ExtractFiles writes === FILE: / === END FILE === blocks from agent output
// into the given directory. If no file blocks are found, the raw output is
// written as main.go.
func ExtractFiles(dir, output string) error {
	matches := fileBlockRe.FindAllStringSubmatch(output, -1)
	if len(matches) == 0 {
		return os.WriteFile(filepath.Join(dir, "main.go"), []byte(output), 0644)
	}
	for _, m := range matches {
		filename := strings.TrimSpace(m[1])
		content := m[2]
		path := filepath.Join(dir, filename)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

// sourceExts lists extensions treated as source code by ExtractFilesFlat.
var sourceExts = map[string]bool{
	".go": true, ".ts": true, ".js": true, ".py": true,
	".tsx": true, ".jsx": true, ".mjs": true, ".cjs": true,
}

// ExtractFilesFlat is like ExtractFiles but enforces a flat layout:
// directory prefixes are stripped from filenames and non-source files
// (e.g. package.json, tsconfig.json) are skipped.
func ExtractFilesFlat(dir, output string) error {
	matches := fileBlockRe.FindAllStringSubmatch(output, -1)
	if len(matches) == 0 {
		return os.WriteFile(filepath.Join(dir, "main.go"), []byte(output), 0644)
	}
	for _, m := range matches {
		filename := filepath.Base(strings.TrimSpace(m[1]))
		ext := filepath.Ext(filename)
		if !sourceExts[ext] {
			continue
		}
		if err := os.WriteFile(filepath.Join(dir, filename), []byte(m[2]), 0644); err != nil {
			return err
		}
	}
	return nil
}

// RunTests executes the language-appropriate test runner in dir and returns
// the parsed coverage percentage and whether tests actually ran.
func RunTests(dir, language string) (coverage float64, ran bool) {
	cmd, args := testCommand(language)
	if cmd == "" {
		return -1, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, cmd, args...)
	c.Dir = dir
	out, _ := c.CombinedOutput()
	return ParseCoverage(string(out)), true
}

// ParseCoverage extracts a Go-style "coverage: NN.N%" value from test output.
func ParseCoverage(output string) float64 {
	m := coverageRe.FindStringSubmatch(output)
	if len(m) < 2 {
		return -1
	}
	v, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return -1
	}
	return v
}

var (
	fileBlockRe = regexp.MustCompile(`(?s)=== FILE: (.+?) ===\n(.+?)=== END FILE ===`)
	coverageRe  = regexp.MustCompile(`coverage:\s+([\d.]+)%`)
)

func testCommand(language string) (string, []string) {
	switch language {
	case "go":
		return "go", []string{"test", "-cover", "."}
	case "ts":
		return "node", []string{"--test"}
	case "py":
		return "python", []string{"-m", "pytest", "-q"}
	default:
		return "", nil
	}
}
