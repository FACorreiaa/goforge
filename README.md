# 🔨 GoForge

A CLI tool to scaffold production-ready Go web applications with an opinionated tech stack.

## 🚀 Stack

Generated projects include:

| Category | Technology |
|----------|------------|
| **Router** | [Chi](https://go-chi.io/) |
| **Templates** | [Templ](https://templ.guide/) |
| **Frontend** | [HTMX](https://htmx.org/) + [Tailwind CSS](https://tailwindcss.com/) + [DaisyUI](https://daisyui.com/) |
| **Database** | PostgreSQL with [pgxpool](https://github.com/jackc/pgx) |
| **Migrations** | [Goose](https://github.com/pressly/goose) |
| **Live Reload** | [Air](https://github.com/air-verse/air) |
| **Releases** | [GoReleaser](https://goreleaser.com/) |

**No Node.js required!** Tailwind uses the standalone CLI.

## 📦 Installation

### From Source

```bash
go install github.com/FACorreiaa/goforge@latest
```

### From Binary

Download from the [Releases](https://github.com/FACorreiaa/goforge/releases) page.

## 🎯 Usage

### Create a New Project

```bash
# With arguments
goforge new my-app github.com/username/my-app

# Interactive mode
goforge new
```

### Generated Project Structure

```
my-app/
├── cmd/server/          # Entry point
├── internal/
│   ├── config/          # Configuration
│   ├── database/        # Pgx pool + migrations
│   ├── middleware/      # HTTP middleware
│   └── server/          # Router & handlers
├── views/               # Templ templates
├── assets/              # CSS, JS, static files
├── pkg/helpers/         # Utility functions
├── Makefile             # Build commands
├── Dockerfile           # Multi-stage build
├── docker-compose.yml   # Dev stack
└── .goreleaser.yml      # Release config
```

### Start Development

```bash
cd my-app
make setup    # Install Air, Templ, Goose, Tailwind CLI
make dev      # Start with live reload
```

## 🛠 Development

### Building GoForge

```bash
git clone https://github.com/FACorreiaa/goforge.git
cd goforge
go build -o goforge .
```

### Testing

```bash
./goforge new test-app github.com/test/test-app
cd test-app
make setup
make dev
```

## 📝 License

MIT License - see [LICENSE](LICENSE)

---

Built with ❤️ in Go
