package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/FACorreiaa/goforge/internal/generator"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
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

var rootCmd = &cobra.Command{
	Use:   "goforge",
	Short: "Scaffold production-ready Go projects",
	Long: `GoForge - A CLI tool to scaffold production-ready Go projects.

Stack includes: Go + Chi + Templ + HTMX + Tailwind + DaisyUI + Postgres + Pgx + Goose

Example:
  goforge new my-app github.com/username/my-app
  goforge new                                      # Interactive mode`,
	Version: version,
}

var newCmd = &cobra.Command{
	Use:   "new [project-name] [module-path]",
	Short: "Create a new project",
	Long:  `Create a new project with the GoForge stack.`,
	Args:  cobra.MaximumNArgs(2),
	RunE:  runNew,
}

// Flags
var (
	frontendFlag     string
	cssFrameworkFlag string
)

func init() {
	newCmd.Flags().StringVarP(&frontendFlag, "frontend", "f", "", "Frontend stack: htmx, htmx-hyperscript, htmx-alpine")
	newCmd.Flags().StringVarP(&cssFrameworkFlag, "css", "c", "", "CSS framework: daisyui, templui, basecoat")
	rootCmd.AddCommand(newCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runNew(cmd *cobra.Command, args []string) error {
	var projectName, modulePath, frontend, cssFramework string

	// Handle arguments
	if len(args) >= 2 {
		projectName = args[0]
		modulePath = args[1]
	} else if len(args) == 1 {
		projectName = args[0]
		// Prompt for module path only
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Module Path").
					Description("Go module path (e.g., github.com/username/project)").
					Value(&modulePath).
					Validate(validateModulePath),
			),
		)
		if err := form.Run(); err != nil {
			return err
		}
	} else {
		// Full interactive mode
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Project Name").
					Description("Name of your project directory").
					Value(&projectName).
					Validate(validateProjectName),
				huh.NewInput().
					Title("Module Path").
					Description("Go module path (e.g., github.com/username/project)").
					Value(&modulePath).
					Validate(validateModulePath),
			),
		)
		if err := form.Run(); err != nil {
			return err
		}
	}

	// Handle frontend selection
	if frontendFlag != "" {
		frontend = frontendFlag
	} else {
		// Prompt for frontend choice
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Frontend Stack").
					Description("Choose your frontend enhancement library").
					Options(
						huh.NewOption("HTMX only", FrontendHTMX),
						huh.NewOption("HTMX + Hyperscript (_hyperscript)", FrontendHTMXHyperscript),
						huh.NewOption("HTMX + Alpine.js", FrontendHTMXAlpine),
					).
					Value(&frontend),
			),
		)
		if err := form.Run(); err != nil {
			return err
		}
	}

	// Validate frontend choice
	if frontend != FrontendHTMX && frontend != FrontendHTMXHyperscript && frontend != FrontendHTMXAlpine {
		frontend = FrontendHTMX // Default to HTMX only
	}

	// Handle CSS framework selection
	if cssFrameworkFlag != "" {
		cssFramework = cssFrameworkFlag
	} else {
		// Prompt for CSS framework choice
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("CSS Framework").
					Description("Choose your CSS component framework").
					Options(
						huh.NewOption("DaisyUI - Component library with themes", CSSFrameworkDaisyUI),
						huh.NewOption("TemplUI - Go/Templ component library", CSSFrameworkTemplUI),
						huh.NewOption("Basecoat - shadcn/ui-style components", CSSFrameworkBasecoat),
					).
					Value(&cssFramework),
			),
		)
		if err := form.Run(); err != nil {
			return err
		}
	}

	// Validate CSS framework choice
	if cssFramework != CSSFrameworkDaisyUI && cssFramework != CSSFrameworkTemplUI && cssFramework != CSSFrameworkBasecoat {
		cssFramework = CSSFrameworkDaisyUI // Default to DaisyUI
	}

	// Get absolute path
	absPath, err := filepath.Abs(projectName)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if directory exists
	if _, err := os.Stat(absPath); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", projectName)
	}

	frontendLabel := map[string]string{
		FrontendHTMX:            "HTMX",
		FrontendHTMXHyperscript: "HTMX + Hyperscript",
		FrontendHTMXAlpine:      "HTMX + Alpine.js",
	}[frontend]

	cssLabel := map[string]string{
		CSSFrameworkDaisyUI:  "DaisyUI",
		CSSFrameworkTemplUI:  "TemplUI",
		CSSFrameworkBasecoat: "Basecoat",
	}[cssFramework]

	fmt.Printf("\n🚀 Creating project '%s' with module '%s'...\n", projectName, modulePath)
	fmt.Printf("   Frontend: %s\n", frontendLabel)
	fmt.Printf("   CSS Framework: %s\n\n", cssLabel)

	// Generate the project with options
	opts := generator.Options{
		ProjectName:  projectName,
		ModulePath:   modulePath,
		Frontend:     frontend,
		CSSFramework: cssFramework,
	}
	if err := generator.GenerateWithOptions(opts); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// Success message
	fmt.Println("\n✅ Project created successfully!")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  make setup    # Install tools (Air, Templ, Goose, Tailwind)")
	fmt.Println("  make dev      # Start development server with live reload")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Println("\n📚 See README.md for more commands and documentation.")

	return nil
}

func validateProjectName(s string) error {
	if s == "" {
		return fmt.Errorf("project name is required")
	}
	if strings.ContainsAny(s, " /\\") {
		return fmt.Errorf("project name cannot contain spaces or slashes")
	}
	return nil
}

func validateModulePath(s string) error {
	if s == "" {
		return fmt.Errorf("module path is required")
	}
	if !strings.Contains(s, "/") {
		return fmt.Errorf("module path should contain at least one '/'")
	}
	return nil
}
