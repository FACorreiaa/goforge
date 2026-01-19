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

func init() {
	rootCmd.AddCommand(newCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runNew(cmd *cobra.Command, args []string) error {
	var projectName, modulePath string

	if len(args) >= 2 {
		projectName = args[0]
		modulePath = args[1]
	} else if len(args) == 1 {
		projectName = args[0]
		// Prompt for module path
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Module Path").
					Description("Go module path (e.g., github.com/username/project)").
					Value(&modulePath).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("module path is required")
						}
						if !strings.Contains(s, "/") {
							return fmt.Errorf("module path should contain at least one '/'")
						}
						return nil
					}),
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
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("project name is required")
						}
						if strings.ContainsAny(s, " /\\") {
							return fmt.Errorf("project name cannot contain spaces or slashes")
						}
						return nil
					}),
				huh.NewInput().
					Title("Module Path").
					Description("Go module path (e.g., github.com/username/project)").
					Value(&modulePath).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("module path is required")
						}
						if !strings.Contains(s, "/") {
							return fmt.Errorf("module path should contain at least one '/'")
						}
						return nil
					}),
			),
		)
		if err := form.Run(); err != nil {
			return err
		}
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

	fmt.Printf("\n🚀 Creating project '%s' with module '%s'...\n\n", projectName, modulePath)

	// Generate the project
	if err := generator.Generate(projectName, modulePath); err != nil {
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
