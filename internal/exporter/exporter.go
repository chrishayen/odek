package exporter

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	runepkg "github.com/chrishayen/odek/internal/rune"
)

// Result holds the outcome of exporting a feature.
type Result struct {
	FeatureName string   `json:"feature_name"`
	Version     string   `json:"version"`
	OutputDir   string   `json:"output_dir"`
	Files       []string `json:"files"`
}

// Exporter assembles a composed feature into a standalone library.
type Exporter struct {
	runeStore *runepkg.Store
	language  string
}

func New(runeStore *runepkg.Store, language string) *Exporter {
	return &Exporter{runeStore: runeStore, language: language}
}

// ExportOptions controls export behavior.
type ExportOptions struct {
	IncludeTests bool
}

// Export bundles a feature's hydrated code into a standalone library at distPath/<feature>/.
func (e *Exporter) Export(name, distPath string, opts ExportOptions) (*Result, error) {
	// Look up the top-level rune directly — it may be structural (no signature).
	topRune, err := e.runeStore.Get(name)
	if err != nil {
		return nil, fmt.Errorf("feature %q not found", name)
	}
	version := topRune.Version.String()

	// Check every leaf rune under this feature is hydrated.
	// Use prefix "name." to avoid matching sibling runes (e.g. "foo" shouldn't match "foo_bar").
	allRunes, err := e.runeStore.ListByPrefix(name + ".")
	if err != nil {
		return nil, fmt.Errorf("listing runes for feature %q: %w", name, err)
	}
	// Include the top-level rune itself.
	children := append(allRunes, *topRune)
	allNames := make([]string, len(children))
	for i, r := range children {
		allNames[i] = r.Name
	}
	var unhydrated []string
	for _, r := range children {
		if runepkg.IsLeaf(r.Name, allNames) && !r.Hydrated {
			unhydrated = append(unhydrated, r.Name)
		}
	}
	if len(unhydrated) > 0 {
		return nil, fmt.Errorf("feature %q has unhydrated runes: %s", name, strings.Join(unhydrated, ", "))
	}

	srcDir := filepath.Join(e.runeStore.OutputPath(), name)
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("feature source directory %q does not exist — hydrate first", srcDir)
	}

	outDir := filepath.Join(distPath, name)
	if err := os.RemoveAll(outDir); err != nil {
		return nil, fmt.Errorf("cleaning dist dir: %w", err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("creating dist dir: %w", err)
	}

	ext := langExtension(e.language)
	testSuffix := testFileSuffix(e.language)

	var files []string

	// Copy source files from src/<feature>/ tree.
	err = filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ext) {
			return nil
		}
		if !opts.IncludeTests && strings.HasSuffix(d.Name(), testSuffix) {
			return nil
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(outDir, rel)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return err
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("copying source files: %w", err)
	}

	// Copy external dependency source files (e.g. std/io/write_stdout.ts).
	depFiles, err := e.copyDependencyFiles(name, children, outDir, ext, testSuffix, opts)
	if err != nil {
		return nil, fmt.Errorf("copying dependency files: %w", err)
	}
	files = append(files, depFiles...)

	// Generate wiring file that injects dependencies.
	wiringContent := e.generateWiring(name, children)
	if wiringContent != "" {
		wiringFile := "wiring" + ext
		if err := os.WriteFile(filepath.Join(outDir, wiringFile), []byte(wiringContent), 0644); err != nil {
			return nil, fmt.Errorf("writing wiring: %w", err)
		}
		files = append(files, wiringFile)
	}

	// Generate index file — re-exports wiring if it exists, otherwise raw source.
	indexContent := e.generateIndex(name, children, files)
	indexFile := "index" + ext
	if err := os.WriteFile(filepath.Join(outDir, indexFile), []byte(indexContent), 0644); err != nil {
		return nil, fmt.Errorf("writing index: %w", err)
	}
	files = append(files, indexFile)

	// Generate package.json.
	pkgJSON := generatePackageJSON(name, version)
	if err := os.WriteFile(filepath.Join(outDir, "package.json"), pkgJSON, 0644); err != nil {
		return nil, fmt.Errorf("writing package.json: %w", err)
	}
	files = append(files, "package.json")

	return &Result{
		FeatureName: name,
		Version:     version,
		OutputDir:   outDir,
		Files:       files,
	}, nil
}

func langExtension(lang string) string {
	switch lang {
	case "go":
		return ".go"
	case "py":
		return ".py"
	default:
		return ".ts"
	}
}

func testFileSuffix(lang string) string {
	switch lang {
	case "go":
		return "_test.go"
	case "py":
		return "_test.py"
	default:
		return ".test.ts"
	}
}

// copyDependencyFiles copies source files for external dependencies (runes outside
// the feature's own tree) into the export directory.
func (e *Exporter) copyDependencyFiles(featureName string, runes []runepkg.Rune, outDir, ext, testSuffix string, opts ExportOptions) ([]string, error) {
	outputPath := e.runeStore.OutputPath()
	seen := make(map[string]bool)
	var files []string

	for _, r := range runes {
		for _, ref := range r.Dependencies {
			depPath, _ := runepkg.ParseRef(ref)
			if depPath == "" {
				depPath = ref
			}
			// Skip deps within the same feature tree.
			if depPath == featureName || strings.HasPrefix(depPath, featureName+".") {
				continue
			}
			if seen[depPath] {
				continue
			}
			seen[depPath] = true

			// Find and copy the dep's source file.
			codeDir := e.runeStore.CodeDir(depPath)
			shortName := runepkg.ShortName(depPath)
			srcFile := filepath.Join(codeDir, shortName+ext)
			if _, err := os.Stat(srcFile); os.IsNotExist(err) {
				continue
			}

			// Compute relative path from outputPath for the dest.
			rel, err := filepath.Rel(outputPath, srcFile)
			if err != nil {
				continue
			}
			destPath := filepath.Join(outDir, rel)
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return nil, err
			}
			data, err := os.ReadFile(srcFile)
			if err != nil {
				return nil, err
			}
			if err := os.WriteFile(destPath, data, 0644); err != nil {
				return nil, err
			}
			files = append(files, rel)

			// Also copy test file if requested.
			if opts.IncludeTests {
				testFile := filepath.Join(codeDir, shortName+testSuffix)
				if _, err := os.Stat(testFile); err == nil {
					testRel, _ := filepath.Rel(outputPath, testFile)
					testDest := filepath.Join(outDir, testRel)
					data, _ := os.ReadFile(testFile)
					os.WriteFile(testDest, data, 0644)
					files = append(files, testRel)
				}
			}
		}
	}
	return files, nil
}

// generateWiring creates a wiring file that imports implementations and injects
// dependencies. Returns empty string if no wiring is needed.
func (e *Exporter) generateWiring(featureName string, runes []runepkg.Rune) string {
	if e.language != "ts" {
		return ""
	}

	// Find runes in this feature's tree that have external dependencies.
	type wiringEntry struct {
		rune     runepkg.Rune
		extDeps  []runepkg.Rune // resolved external dependency runes
	}
	var entries []wiringEntry

	for _, r := range runes {
		if len(r.Dependencies) == 0 {
			continue
		}
		// Only wire runes in this feature's tree.
		if r.Name != featureName && !strings.HasPrefix(r.Name, featureName+".") {
			continue
		}
		var extDeps []runepkg.Rune
		for _, ref := range r.Dependencies {
			depPath, _ := runepkg.ParseRef(ref)
			if depPath == "" {
				depPath = ref
			}
			if depPath == featureName || strings.HasPrefix(depPath, featureName+".") {
				continue // internal dep, skip for now
			}
			dep, err := e.runeStore.Get(depPath)
			if err != nil {
				continue
			}
			extDeps = append(extDeps, *dep)
		}
		if len(extDeps) > 0 {
			entries = append(entries, wiringEntry{rune: r, extDeps: extDeps})
		}
	}

	if len(entries) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("// Auto-generated wiring — connects runes to their dependencies.\n\n")

	// Collect all imports.
	importedDeps := make(map[string]bool)
	for _, entry := range entries {
		for _, dep := range entry.extDeps {
			if importedDeps[dep.Name] {
				continue
			}
			importedDeps[dep.Name] = true
			short := runepkg.ShortName(dep.Name)
			// Path from feature root to the dep's source file.
			parts := strings.Split(dep.Name, ".")
			importPath := "./" + filepath.Join(parts[:len(parts)-1]...) + "/" + short
			fmt.Fprintf(&b, "import { %s } from '%s';\n", short, importPath)
		}
	}

	for _, entry := range entries {
		short := runepkg.ShortName(entry.rune.Name)
		// Import path for the rune itself (relative to feature root).
		var runeImportPath string
		if entry.rune.Name == featureName {
			// Top-level rune, file is in the feature root.
			runeImportPath = "./" + short
		} else {
			// Nested rune.
			rel := strings.TrimPrefix(entry.rune.Name, featureName+".")
			parts := strings.Split(rel, ".")
			if len(parts) > 1 {
				runeImportPath = "./" + filepath.Join(parts[:len(parts)-1]...) + "/" + short
			} else {
				runeImportPath = "./" + short
			}
		}
		fmt.Fprintf(&b, "import { %s as _%s } from '%s';\n", short, short, runeImportPath)
	}

	b.WriteString("\n")

	// Generate wired exports.
	for _, entry := range entries {
		short := runepkg.ShortName(entry.rune.Name)
		depArgs := make([]string, len(entry.extDeps))
		for i, dep := range entry.extDeps {
			depArgs[i] = runepkg.ShortName(dep.Name)
		}
		fmt.Fprintf(&b, "export function %s(...args: Parameters<typeof _%s>) {\n", short, short)
		fmt.Fprintf(&b, "  return _%s(%s, ...args);\n", short, strings.Join(depArgs, ", "))
		fmt.Fprintf(&b, "}\n\n")
	}

	return b.String()
}

// generateIndex creates an index file for the export.
func (e *Exporter) generateIndex(featureName string, runes []runepkg.Rune, files []string) string {
	if e.language != "ts" {
		return ""
	}

	// If we generated wiring, re-export from wiring.
	hasWiring := false
	for _, f := range files {
		if f == "wiring.ts" {
			hasWiring = true
			break
		}
	}

	var b strings.Builder
	if hasWiring {
		b.WriteString("export * from './wiring';\n")
	} else {
		// Re-export top-level source files directly.
		ext := langExtension(e.language)
		for _, f := range files {
			if strings.Contains(f, string(filepath.Separator)) {
				continue
			}
			if !strings.HasSuffix(f, ext) {
				continue
			}
			if f == "index"+ext || f == "wiring"+ext {
				continue
			}
			module := strings.TrimSuffix(f, ext)
			fmt.Fprintf(&b, "export * from './%s';\n", module)
		}
	}
	return b.String()
}

func generatePackageJSON(name, version string) []byte {
	pkg := map[string]any{
		"name":    name,
		"version": version,
		"type":    "module",
		"main":    "index.ts",
	}
	data, _ := json.MarshalIndent(pkg, "", "  ")
	return append(data, '\n')
}
