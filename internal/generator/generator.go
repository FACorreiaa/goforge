package generator

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed all:templates
var templateFS embed.FS

const (
	// Placeholder module used in the golden templates
	placeholderModule = "github.com/goforge/scaffold"

	// Frontend placeholders in templates
	placeholderFrontendScripts = "<!-- FRONTEND_SCRIPTS -->"
)

// Frontend constants
const (
	FrontendHTMX            = "htmx"
	FrontendHTMXHyperscript = "htmx-hyperscript"
	FrontendHTMXAlpine      = "htmx-alpine"
)

// Options for project generation
type Options struct {
	ProjectName string
	ModulePath  string
	Frontend    string // htmx, htmx-hyperscript, htmx-alpine
}

// Generate creates a new project from the embedded templates (backward compatible)
func Generate(projectName string, newModule string) error {
	return GenerateWithOptions(Options{
		ProjectName: projectName,
		ModulePath:  newModule,
		Frontend:    FrontendHTMX,
	})
}

// GenerateWithOptions creates a new project with custom options
func GenerateWithOptions(opts Options) error {
	// Create the project directory
	if err := os.MkdirAll(opts.ProjectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Get frontend scripts based on selection
	frontendScripts := getFrontendScripts(opts.Frontend)

	// Walk through the embedded templates
	err := fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path (remove "templates/" prefix)
		relPath := strings.TrimPrefix(path, "templates/")
		if relPath == "" || relPath == "templates" {
			return nil // Skip the root folder itself
		}

		// Determine target path on user's disk
		targetPath := filepath.Join(opts.ProjectName, relPath)

		// Handle Directories
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Handle Files
		data, err := templateFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", path, err)
		}

		// Perform content replacement for text files
		content := string(data)

		// Skip binary files (basic check)
		if !isBinaryFile(path) {
			// Replace module path
			content = strings.ReplaceAll(content, placeholderModule, opts.ModulePath)
			// Replace frontend scripts placeholder
			content = strings.ReplaceAll(content, placeholderFrontendScripts, frontendScripts)
		}

		// Handle .tmpl extension (strip it from the target)
		if strings.HasSuffix(targetPath, ".tmpl") {
			targetPath = strings.TrimSuffix(targetPath, ".tmpl")
		}

		// Write to disk
		if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		fmt.Printf("  ✓ %s\n", strings.TrimPrefix(targetPath, opts.ProjectName+"/"))
		return nil
	})

	return err
}

// getFrontendScripts returns the appropriate script tags for the selected frontend
func getFrontendScripts(frontend string) string {
	htmxScript := `<script src="https://unpkg.com/htmx.org@2.0.4" integrity="sha384-HGfztofotfshcF7+8n44JQL2oJmowVChPTg48S+jvZoztPfvwD79OC/LTtG6dMp+" crossorigin="anonymous"></script>`

	switch frontend {
	case FrontendHTMXHyperscript:
		return htmxScript + `
			<!-- Hyperscript - _hyperscript -->
			<script src="https://unpkg.com/hyperscript.org@0.9.14"></script>`
	case FrontendHTMXAlpine:
		return htmxScript + `
			<!-- Alpine.js -->
			<script defer src="https://unpkg.com/alpinejs@3.14.8/dist/cdn.min.js"></script>`
	default: // FrontendHTMX
		return htmxScript
	}
}

// isBinaryFile checks if a file is likely binary based on extension
func isBinaryFile(path string) bool {
	binaryExtensions := []string{
		".png", ".jpg", ".jpeg", ".gif", ".ico", ".webp",
		".woff", ".woff2", ".ttf", ".eot",
		".zip", ".tar", ".gz",
		".mjs", // Skip mjs files as they're large
	}

	ext := strings.ToLower(filepath.Ext(path))
	for _, binExt := range binaryExtensions {
		if ext == binExt {
			return true
		}
	}
	return false
}
