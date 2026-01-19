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
)

// Generate creates a new project from the embedded templates
func Generate(projectName string, newModule string) error {
	// Create the project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

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
		targetPath := filepath.Join(projectName, relPath)

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
			content = strings.ReplaceAll(content, placeholderModule, newModule)
		}

		// Handle .tmpl extension (strip it from the target)
		if strings.HasSuffix(targetPath, ".tmpl") {
			targetPath = strings.TrimSuffix(targetPath, ".tmpl")
		}

		// Write to disk
		if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		fmt.Printf("  ✓ %s\n", strings.TrimPrefix(targetPath, projectName+"/"))
		return nil
	})

	return err
}

// isBinaryFile checks if a file is likely binary based on extension
func isBinaryFile(path string) bool {
	binaryExtensions := []string{
		".png", ".jpg", ".jpeg", ".gif", ".ico", ".webp",
		".woff", ".woff2", ".ttf", ".eot",
		".zip", ".tar", ".gz",
	}

	ext := strings.ToLower(filepath.Ext(path))
	for _, binExt := range binaryExtensions {
		if ext == binExt {
			return true
		}
	}
	return false
}
