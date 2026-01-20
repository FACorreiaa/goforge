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

	// Placeholders
	placeholderFrontendScripts = "<!-- FRONTEND_SCRIPTS -->"
	placeholderSetupCommand    = "<!-- SETUP_COMMAND -->"
	placeholderCiSetupCommand  = "<!-- CI_SETUP_COMMAND -->"
	placeholderDevCommand      = "<!-- DEV_COMMAND -->"
	placeholderCssWatchCmd     = "<!-- CSS_WATCH_COMMAND -->"
	placeholderCssBuildCmd     = "<!-- CSS_BUILD_COMMAND -->"
	placeholderAirBuildCmd     = "<!-- AIR_BUILD_CMD -->"
	placeholderTailwindPlugin  = "<!-- TAILWIND_PLUGIN -->"
	placeholderDaisyuiConfig   = "<!-- DAISYUI_CONFIG -->"
	placeholderDockerSetupRun  = "<!-- DOCKER_SETUP_RUN -->"
	placeholderDockerBuildCss  = "<!-- DOCKER_BUILD_CSS -->"
)

// Frontend options
const (
	FrontendHTMX            = "htmx"
	FrontendHTMXHyperscript = "htmx-hyperscript"
	FrontendHTMXAlpine      = "htmx-alpine"
)

// CSS Framework options
const (
	CSSFrameworkDaisyUI  = "daisyui"
	CSSFrameworkTemplUI  = "templui"
	CSSFrameworkBasecoat = "basecoat"
)

// Options for project generation
type Options struct {
	ProjectName  string
	ModulePath   string
	Frontend     string
	CSSFramework string
	IncludeDB    bool
}

// Generate creates a new project from the embedded templates (backward compatible)
func Generate(projectName string, newModule string) error {
	return GenerateWithOptions(Options{
		ProjectName:  projectName,
		ModulePath:   newModule,
		Frontend:     FrontendHTMX,
		CSSFramework: CSSFrameworkDaisyUI,
		IncludeDB:    true, // Default to true for backward compatibility
	})
}

// GenerateWithOptions creates a new project with custom options
func GenerateWithOptions(opts Options) error {
	// Create the project directory
	if err := os.MkdirAll(opts.ProjectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Prepare replacements
	replacements := getReplacements(opts)

	// Walk through the embedded templates
	return fs.WalkDir(templateFS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path (remove "templates/" prefix)
		relPath := strings.TrimPrefix(path, "templates/")
		if relPath == "" || relPath == "templates" {
			return nil // Skip replace root
		}

		// Skip database directory if IncludeDB is false
		if !opts.IncludeDB && (strings.HasPrefix(relPath, "internal/database") ||
			strings.HasPrefix(relPath, "docker-compose.yml") && false) {
			// Check logic for docker-compose: if DB is optional, do we skip file?
			// docker-compose usually provides DB. If no DB, we might want it for other things?
			// Current plan says "Wrap db service in <!-- IF DB -->". So don't skip file.
			// Only skip internal/database folder.
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if !opts.IncludeDB && strings.HasPrefix(relPath, "internal/database") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
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

		if !isBinaryFile(path) {
			// Process conditional blocks first
			content = processConditionalBlocks(content, opts)

			// Replace module path
			content = strings.ReplaceAll(content, placeholderModule, opts.ModulePath)

			// Replace all other placeholders
			for k, v := range replacements {
				content = strings.ReplaceAll(content, k, v)
			}
		}

		// Handle .tmpl extension (strip it from the target)
		targetPath = strings.TrimSuffix(targetPath, ".tmpl")

		// Write to disk
		if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		fmt.Printf("  ✓ %s\n", strings.TrimPrefix(targetPath, opts.ProjectName+"/"))
		return nil
	})
}

// processConditionalBlocks removes content between <!-- IF DB --> and <!-- ENDIF --> if condition is false
// Also supports <!-- IF NOT DB -->
func processConditionalBlocks(content string, opts Options) string {
	// Process DB blocks
	if opts.IncludeDB {
		// Keep content, remove tags
		content = removeTags(content, "<!-- IF DB -->", "<!-- ENDIF -->")
		// Remove NOT DB content
		content = removeBlock(content, "<!-- IF NOT DB -->", "<!-- ENDIF -->")
	} else {
		// Remove DB content
		content = removeBlock(content, "<!-- IF DB -->", "<!-- ENDIF -->")
		// Keep NOT DB content, remove tags
		content = removeTags(content, "<!-- IF NOT DB -->", "<!-- ENDIF -->")
	}
	// Note: Simple regex or string manipulation.
	// Nested blocks not supported with simple logic, but sufficient for this use case.

	return content
}

func removeBlock(content, startTag, endTag string) string {
	for {
		startIndex := strings.Index(content, startTag)
		if startIndex == -1 {
			break
		}
		endIndex := strings.Index(content, endTag)
		if endIndex == -1 {
			break
		}
		// Include tags in removal
		// Find end of endTag
		endTagEnd := endIndex + len(endTag)
		if endIndex > startIndex {
			content = content[:startIndex] + content[endTagEnd:]
		} else {
			// Malformed or nested incorrectly, just break to avoid infinite loop
			break
		}
	}
	return content
}

// removeTags removes the tags but keeps the content inside
func removeTags(content, startTag, endTag string) string {
	content = strings.ReplaceAll(content, startTag, "")
	content = strings.ReplaceAll(content, endTag, "")
	return content
}

func getReplacements(opts Options) map[string]string {
	replacements := make(map[string]string)

	// Frontend JS Downloads
	jsDownloads := `
	@echo "📥 Downloading Frontend Libraries..."
	@curl -sL -o assets/js/htmx.min.js https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js`

	dockerJsDownloads := `RUN curl -sL -o assets/js/htmx.min.js https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js`

	if opts.Frontend == FrontendHTMXHyperscript {
		jsDownloads += `
	@curl -sL -o assets/js/hyperscript.min.js https://unpkg.com/hyperscript.org@0.9.14`
		dockerJsDownloads += ` && \
    curl -sL -o assets/js/hyperscript.min.js https://unpkg.com/hyperscript.org@0.9.14`
	} else if opts.Frontend == FrontendHTMXAlpine {
		jsDownloads += `
	@curl -sL -o assets/js/alpinejs.min.js https://unpkg.com/alpinejs@3.14.8/dist/cdn.min.js`
		dockerJsDownloads += ` && \
    curl -sL -o assets/js/alpinejs.min.js https://unpkg.com/alpinejs@3.14.8/dist/cdn.min.js`
	}

	if opts.CSSFramework == CSSFrameworkBasecoat {
		jsDownloads += `
	@curl -sL -o assets/js/basecoat.min.js https://cdn.jsdelivr.net/npm/basecoat-css@latest/dist/basecoat.min.js`
		dockerJsDownloads += ` && \
    curl -sL -o assets/js/basecoat.min.js https://cdn.jsdelivr.net/npm/basecoat-css@latest/dist/basecoat.min.js`
	}

	// Frontend Scripts (Local Links)
	replacements[placeholderFrontendScripts] = getFrontendScripts(opts)

	// Default Setup (DaisyUI)
	setupCmd := fmt.Sprintf(`@echo "📥 Installing Tailwind CSS + DaisyUI..."
	@mkdir -p assets/css assets/js
	@cd assets && curl -sL daisyui.com/fast | bash
	@mv assets/input.css assets/css/input.css 2>/dev/null || true
	@mv assets/output.css assets/css/output.css 2>/dev/null || true
	@mv assets/daisyui.mjs assets/js/daisyui.mjs 2>/dev/null || true
	@mv assets/daisyui-theme.mjs assets/js/daisyui-theme.mjs 2>/dev/null || true
	@mv assets/tailwindcss ./tailwindcss 2>/dev/null || true
	@if [ -f assets/css/input.css ]; then \
		sed -i.bak 's|./daisyui.mjs|../js/daisyui.mjs|g' assets/css/input.css && rm assets/css/input.css.bak; \
	fi%s`, jsDownloads)

	dockerSetupRun := fmt.Sprintf(`RUN cd assets && curl -sL daisyui.com/fast | bash
# Organize assets
RUN mkdir -p assets/css assets/js && \
    mv assets/input.css assets/css/ && \
    mv assets/output.css assets/css/ && \
    mv assets/daisyui.mjs assets/js/ && \
    mv assets/daisyui-theme.mjs assets/js/ && \
    mv assets/tailwindcss .
# Fix imports
RUN sed -i 's|./daisyui.mjs|../js/daisyui.mjs|g' assets/css/input.css
# Download JS
%s`, dockerJsDownloads)

	devCmd := `@make -j2 dev-air dev-tailwind`
	cssWatchCmd := `@./tailwindcss -i assets/css/input.css -o assets/css/output.css --watch`
	cssBuildCmd := `@./tailwindcss -i assets/css/input.css -o assets/css/output.css`
	airBuildCmd := `templ generate && go build -o ./tmp/main ./cmd/server`

	// Default Docker Build CSS (DaisyUI)
	dockerBuildCss := `RUN ./tailwindcss -i assets/css/input.css -o assets/css/output.css --minify`

	tailwindPlugin := ""
	daisyuiConfig := ""

	// Overrides for Basecoat
	if opts.CSSFramework == CSSFrameworkBasecoat {
		setupCmd = fmt.Sprintf(`@echo "📥 Installing Tailwind CSS + Basecoat..."
	@mkdir -p assets/css assets/js
	@cd assets && curl -sL daisyui.com/fast | bash
	@mv assets/tailwindcss ./tailwindcss 2>/dev/null || true
	@echo "Downloading Basecoat CSS..."
	@curl -sL -o assets/css/basecoat.min.css https://cdn.jsdelivr.net/npm/basecoat-css@latest/dist/basecoat.min.css
	@echo "Creating assets/css/index.css..."
	@echo '@import "tailwindcss"; @import "./basecoat.min.css";' > assets/css/index.css
	@echo "Cleaning up DaisyUI files..."
	@rm assets/input.css assets/output.css assets/daisyui.mjs assets/daisyui-theme.mjs 2>/dev/null || true%s`, jsDownloads)

		dockerSetupRun = fmt.Sprintf(`RUN cd assets && curl -sL daisyui.com/fast | bash
RUN mkdir -p assets/css assets/js && \
    mv assets/tailwindcss .
RUN curl -sL -o assets/css/basecoat.min.css https://cdn.jsdelivr.net/npm/basecoat-css@latest/dist/basecoat.min.css
RUN echo '@import "tailwindcss"; @import "./basecoat.min.css";' > assets/css/index.css
RUN rm assets/input.css assets/output.css assets/daisyui* 2>/dev/null || true
# Download JS
%s`, dockerJsDownloads)

		devCmd = `@make dev-air` // Air handles build
		cssWatchCmd = `@echo "CSS watching handled by Air"`
		cssBuildCmd = `@./tailwindcss -i assets/css/index.css -o assets/css/output.css`

		airBuildCmd = `templ generate && ./tailwindcss -i ./assets/css/index.css -o ./assets/css/output.css --minify && go build -o ./tmp/main ./cmd/server`

		dockerBuildCss = `RUN ./tailwindcss -i assets/css/index.css -o assets/css/output.css --minify`
	}

	// Assign values
	replacements[placeholderSetupCommand] = setupCmd
	replacements[placeholderCiSetupCommand] = setupCmd
	replacements[placeholderDockerSetupRun] = dockerSetupRun

	replacements[placeholderDevCommand] = devCmd
	replacements[placeholderCssWatchCmd] = cssWatchCmd
	replacements[placeholderCssBuildCmd] = cssBuildCmd
	replacements[placeholderAirBuildCmd] = airBuildCmd
	replacements[placeholderDockerBuildCss] = dockerBuildCss

	replacements[placeholderTailwindPlugin] = tailwindPlugin
	replacements[placeholderDaisyuiConfig] = daisyuiConfig

	return replacements
}

// getFrontendScripts returns the appropriate script tags for the selected frontend (Local Files)
func getFrontendScripts(opts Options) string {
	htmxScript := `<script src="/assets/js/htmx.min.js"></script>`
	scripts := htmxScript

	switch opts.Frontend {
	case FrontendHTMXHyperscript:
		scripts += `
			<!-- Hyperscript -->
			<script src="/assets/js/hyperscript.min.js"></script>`
	case FrontendHTMXAlpine:
		scripts += `
			<!-- Alpine.js -->
			<script defer src="/assets/js/alpinejs.min.js"></script>`
	}

	if opts.CSSFramework == CSSFrameworkBasecoat {
		scripts += `
			<!-- Basecoat JS -->
			<script defer src="/assets/js/basecoat.min.js"></script>`
	}

	return scripts
}

// isBinaryFile checks if a file is likely binary based on extension
func isBinaryFile(path string) bool {
	binaryExtensions := []string{
		".png", ".jpg", ".jpeg", ".gif", ".ico", ".webp",
		".woff", ".woff2", ".ttf", ".eot",
		".zip", ".tar", ".gz",
		".mjs",
	}

	ext := strings.ToLower(filepath.Ext(path))
	for _, binExt := range binaryExtensions {
		if ext == binExt {
			return true
		}
	}
	return false
}
