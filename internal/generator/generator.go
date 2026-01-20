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

	// CSS Framework placeholders in templates
	placeholderCSSInputImport  = "<!-- CSS_INPUT_IMPORT -->"
	placeholderCSSEmbedPath    = "<!-- CSS_EMBED_PATH -->"
	placeholderTailwindPlugin  = "<!-- TAILWIND_PLUGIN -->"
	placeholderDaisyUIConfig   = "<!-- DAISYUI_CONFIG -->"
	placeholderSetupCommand    = "<!-- SETUP_COMMAND -->"
	placeholderCSSBuildCommand = "<!-- CSS_BUILD_COMMAND -->"
)

// Frontend constants
const (
	FrontendHTMX            = "htmx"
	FrontendHTMXHyperscript = "htmx-hyperscript"
	FrontendHTMXAlpine      = "htmx-alpine"
)

// CSS Framework constants
const (
	CSSFrameworkDaisyUI = "daisyui"
	CSSFrameworkTemplUI = "templui"
	CSSFrameworkBasecoat = "basecoat"
)

// Options for project generation
type Options struct {
	ProjectName  string
	ModulePath   string
	Frontend     string // htmx, htmx-hyperscript, htmx-alpine
	CSSFramework string // daisyui, templui, basecoat
}

// Generate creates a new project from the embedded templates (backward compatible)
func Generate(projectName string, newModule string) error {
	return GenerateWithOptions(Options{
		ProjectName:  projectName,
		ModulePath:   newModule,
		Frontend:     FrontendHTMX,
		CSSFramework: CSSFrameworkDaisyUI,
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

	// Get CSS framework configuration
	cssConfig := getCSSFrameworkConfig(opts.CSSFramework)

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
			// Replace CSS framework placeholders
			content = strings.ReplaceAll(content, placeholderCSSInputImport, cssConfig.InputImport)
			content = strings.ReplaceAll(content, placeholderCSSEmbedPath, cssConfig.EmbedPath)
			content = strings.ReplaceAll(content, placeholderTailwindPlugin, cssConfig.TailwindPlugin)
			content = strings.ReplaceAll(content, placeholderDaisyUIConfig, cssConfig.DaisyUIConfig)
			content = strings.ReplaceAll(content, placeholderSetupCommand, cssConfig.SetupCommand)
			content = strings.ReplaceAll(content, placeholderCSSBuildCommand, cssConfig.BuildCommand)
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

// CSSFrameworkConfig holds configuration for different CSS frameworks
type CSSFrameworkConfig struct {
	InputImport    string // Content for input.css @import/@plugin
	EmbedPath      string // Path for embed.FS in efs.go
	TailwindPlugin string // Plugin config for tailwind.config.js
	DaisyUIConfig  string // DaisyUI config for tailwind.config.js
	SetupCommand   string // Setup command for Makefile
	BuildCommand   string // Build command for CSS
}

// getCSSFrameworkConfig returns configuration based on CSS framework choice
func getCSSFrameworkConfig(framework string) CSSFrameworkConfig {
	switch framework {
	case CSSFrameworkDaisyUI:
		return CSSFrameworkConfig{
			InputImport: `@import "tailwindcss";

@source not "./tailwindcss";
@source not "./daisyui{,*}.mjs";

@plugin "../js/daisyui.mjs";`,
			EmbedPath:      "css/output.css js/*.mjs static/*",
			TailwindPlugin: `require('./assets/js/daisyui.mjs')`,
			DaisyUIConfig: `,
	daisyui: {
		themes: ["light", "dark"],
		darkTheme: "dark",
		base: true,
		styled: true,
		utils: true,
	}`,
			SetupCommand: `@echo "📥 Installing Tailwind CSS + DaisyUI..."
	@cd assets && curl -sL daisyui.com/fast | bash`,
			BuildCommand: `@cd assets && ./tailwindcss -i css/input.css -o css/output.css`,
		}

	case CSSFrameworkTemplUI:
		return CSSFrameworkConfig{
			InputImport: `@import "tailwindcss";

/* TemplUI Base Styles */
@theme {
	--color-background: oklch(100% 0 0);
	--color-foreground: oklch(10% 0 0);
	--color-primary: oklch(50% 0.2 250);
	--color-secondary: oklch(70% 0.15 200);
}

@theme dark {
	--color-background: oklch(10% 0 0);
	--color-foreground: oklch(95% 0 0);
}`,
			EmbedPath:      "css/output.css static/*",
			TailwindPlugin: "",
			DaisyUIConfig:  "",
			SetupCommand: `@echo "📥 Installing Tailwind CSS..."
	@cd assets && curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$$(uname -s | tr '[:upper:]' '[:lower:]')-$$(uname -m | sed 's/x86_64/x64/;s/aarch64/arm64/') -Lo tailwindcss && chmod +x tailwindcss
	@echo "📦 Installing TemplUI..."
	@go install github.com/templui/templui/cmd/templui@latest`,
			BuildCommand: `@cd assets && ./tailwindcss -i css/input.css -o css/output.css`,
		}

	case CSSFrameworkBasecoat:
		return CSSFrameworkConfig{
			InputImport: `@import "tailwindcss";
@import "basecoat-css";`,
			EmbedPath:      "css/output.css static/*",
			TailwindPlugin: "",
			DaisyUIConfig:  "",
			SetupCommand: `@echo "📥 Installing Tailwind CSS..."
	@cd assets && curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$$(uname -s | tr '[:upper:]' '[:lower:]')-$$(uname -m | sed 's/x86_64/x64/;s/aarch64/arm64/') -Lo tailwindcss && chmod +x tailwindcss
	@echo "📦 Installing Basecoat..."
	@npm install basecoat-css`,
			BuildCommand: `@cd assets && ./tailwindcss -i css/input.css -o css/output.css`,
		}

	default:
		// Default to DaisyUI
		return getCSSFrameworkConfig(CSSFrameworkDaisyUI)
	}
}
