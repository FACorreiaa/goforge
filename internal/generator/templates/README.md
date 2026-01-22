# 🚀 GoForge Project

A production-ready Go web application built with:

- **[Chi](https://github.com/go-chi/chi)** - Lightweight, composable router
- **[Templ](https://templ.guide)** - Type-safe HTML templates
- **[HTMX](https://htmx.org)** - High-powered hypermedia
- **[Tailwind CSS](https://tailwindcss.com)** - Utility-first CSS (standalone, no Node.js!)
- **[DaisyUI](https://daisyui.com)** - Beautiful component library
- **[PostgreSQL](https://postgresql.org)** - Robust database with pgxpool
- **[Goose](https://github.com/pressly/goose)** - Database migrations

## 📋 Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Make

## 🚀 Quick Start

```bash
# 1. Install tools (Air, Templ, Goose, Tailwind)
make setup

# 2. Copy environment file
cp .env.example .env
# Edit .env with your database credentials

# 3. Run database migrations
make db-up

# 4. Start development server (with live reload)
make dev
```

Open [http://localhost:8080](http://localhost:8080) in your browser.

## 📁 Project Structure

```
.
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database connection & migrations
│   ├── middleware/       # HTTP middleware
│   └── server/           # HTTP server & routes
├── views/                # Templ templates
│   ├── layouts/          # Base HTML layouts
│   ├── pages/            # Page templates
│   └── components/       # Reusable components
├── assets/               # Static assets (CSS, JS, images)
│   ├── css/              # Tailwind input
│   ├── dist/             # Generated CSS
│   └── static/           # PWA manifest, icons
├── pkg/helpers/          # Utility functions
├── Makefile              # Build commands
├── Dockerfile            # Production Docker image
└── docker-compose.yml    # Local development stack
```

## 🛠 Available Commands

```bash
# Development
make dev              # Start with live reload (Air)
make dev-templ        # Start with Templ proxy (auto browser refresh)
make templ            # Generate Templ templates

# Building
make build            # Build production binary
make run              # Build and run

# Database
make db-up            # Run migrations
make db-down          # Rollback last migration
make db-status        # Show migration status
make db-create        # Create new migration

# Docker
make docker-build     # Build Docker image
make docker-run       # Run Docker container
make docker-compose-up    # Start with docker-compose

# Quality
make test             # Run tests
make lint             # Run golangci-lint
make fmt              # Format code (templ + gofumpt)
make test-coverage    # Run tests with coverage

# Utilities
make clean            # Remove build artifacts
make help             # Show all commands
```

<!-- IF HOOKS -->
## 🪝 Git Hooks

This project includes pre-commit hooks for code quality:

| Hook | Actions |
|------|---------|
| **pre-commit** | `templ fmt` → `templ generate` → `gofumpt` → `golangci-lint` |

The hooks are automatically installed during `make setup`. To reinstall manually:

```bash
make setup-hooks
```

To bypass hooks temporarily (not recommended):

```bash
git commit --no-verify
```
<!-- /IF HOOKS -->

## 🔧 Configuration

Environment variables (see `.env.example`):

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `GO_ENV` | Environment (development/production) | `development` |
| `DATABASE_URL` | PostgreSQL connection string | - |

## 🎨 Styling

This project uses **Tailwind CSS Standalone** - no Node.js required!

- Input CSS: `assets/css/input.css`
- Output CSS: `assets/dist/styles.css` (generated)
- Config: `tailwind.config.js`

DaisyUI is loaded via CDN for simplicity.

## 📱 PWA Support

The app is installable as a Progressive Web App:

- Manifest: `assets/static/manifest.json`
- Service Worker: `assets/static/sw.js`

## 🚢 Deployment

### Docker

```bash
# Build and run with Docker Compose
docker-compose up -d

# Or build manually
docker build -t myapp .
docker run -p 8080:8080 --env-file .env myapp
```

### GoReleaser

```bash
# Create a release
goreleaser release --clean
```

<!-- IF DEPLOY_HETZNER -->
### Hetzner + Caddy

This project includes deployment configuration for Hetzner VPS with Caddy reverse proxy.

**First-time setup:**
```bash
# 1. Add to .env:
#    DEPLOY_HOST=root@your-server-ip
#    DEPLOY_PATH=/opt/myapp

# 2. Setup the server (installs Caddy, configures firewall)
make deploy-setup

# 3. Edit deploy/Caddyfile and replace YOUR_DOMAIN

# 4. Upload .env to server
scp .env root@your-server-ip:/opt/myapp/.env
```

**Deploy updates:**
```bash
make deploy
```

The deploy script will:
1. Build the production binary for Linux
2. Upload binary and assets via SSH
3. Restart the systemd service
<!-- /IF DEPLOY_HETZNER -->

## 📚 Resources

- [Chi Documentation](https://go-chi.io/)
- [Templ Guide](https://templ.guide/)
- [HTMX Documentation](https://htmx.org/docs/)
- [DaisyUI Components](https://daisyui.com/components/)

---

Built with ❤️ using [GoForge](https://github.com/fernando-idwell/goforge)
