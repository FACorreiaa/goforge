package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir := t.TempDir()
	projectName := filepath.Join(tmpDir, "test-app")
	moduleName := "github.com/test/test-app"

	// Generate the project
	err := Generate(projectName, moduleName)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Define expected files
	expectedFiles := []string{
		"go.mod",
		"Makefile",
		"README.md",
		".air.toml",
		".env.example",
		".gitignore",
		".golangci.yml",
		".goreleaser.yml",
		"Dockerfile",
		"docker-compose.yml",
		"tailwind.config.js",
		"cmd/server/main.go",
		"internal/server/server.go",
		"internal/server/routes.go",
		"internal/config/config.go",
		"internal/database/database.go",
		"internal/database/migrations/00001_init.sql",
		"internal/middleware/middleware.go",
		"internal/middleware/logger.go",
		"internal/middleware/secure.go",
		"internal/middleware/ratelimit.go",
		"pkg/helpers/response.go",
		"pkg/helpers/validate.go",
		"views/layouts/base.templ",
		"views/pages/index.templ",
		"views/components/navbar.templ",
		"views/components/footer.templ",
		"assets/efs.go",
		"assets/input.css",
		"assets/output.css",
		"assets/daisyui.mjs",
		"assets/daisyui-theme.mjs",
		"assets/static/manifest.json",
		"assets/static/sw.js",
	}

	// Check that all expected files exist
	for _, file := range expectedFiles {
		path := filepath.Join(projectName, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file missing: %s", file)
		}
	}

	// Check go.mod has correct module name
	goModPath := filepath.Join(projectName, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}
	if !strings.Contains(string(content), moduleName) {
		t.Errorf("go.mod does not contain module name %s", moduleName)
	}

	// Check imports are correctly replaced
	mainGoPath := filepath.Join(projectName, "cmd/server/main.go")
	mainContent, err := os.ReadFile(mainGoPath)
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}
	if !strings.Contains(string(mainContent), moduleName) {
		t.Errorf("main.go imports do not use module name %s", moduleName)
	}
	if strings.Contains(string(mainContent), "github.com/goforge/scaffold") {
		t.Error("main.go still contains placeholder module name")
	}

	// Check Templ files have correct imports
	indexTemplPath := filepath.Join(projectName, "views/pages/index.templ")
	indexContent, err := os.ReadFile(indexTemplPath)
	if err != nil {
		t.Fatalf("Failed to read index.templ: %v", err)
	}
	if !strings.Contains(string(indexContent), moduleName) {
		t.Errorf("index.templ imports do not use module name %s", moduleName)
	}
}

func TestGenerateProjectAlreadyExists(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	projectName := filepath.Join(tmpDir, "existing-project")

	// Create the directory first
	err := os.MkdirAll(projectName, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Try to generate - should still work (we allow generating into existing empty dirs)
	err = Generate(projectName, "github.com/test/existing-project")
	if err != nil {
		t.Fatalf("Generate into existing empty directory failed: %v", err)
	}
}

func TestGenerateFileCount(t *testing.T) {
	tmpDir := t.TempDir()
	projectName := filepath.Join(tmpDir, "count-test")
	moduleName := "github.com/test/count-test"

	err := Generate(projectName, moduleName)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Count files (excluding directories)
	var fileCount int
	err = filepath.Walk(projectName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileCount++
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk project: %v", err)
	}

	// We expect at least 30 files
	minExpectedFiles := 30
	if fileCount < minExpectedFiles {
		t.Errorf("Expected at least %d files, got %d", minExpectedFiles, fileCount)
	}
}
