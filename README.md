<div align="center">

# ⚡ GoMigrate

**Beautiful, interactive database migrations for Go developers**

[![Go Reference](https://pkg.go.dev/badge/github.com/taqnihub/gomigrate.svg)](https://pkg.go.dev/github.com/taqnihub/gomigrate)
[![Go Report Card](https://goreportcard.com/badge/github.com/taqnihub/gomigrate)](https://goreportcard.com/report/github.com/taqnihub/gomigrate)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8.svg)](https://golang.org/)

A friendly wrapper around [golang-migrate](https://github.com/golang-migrate/migrate) with short commands, YAML config, and a polished terminal UI.

```bash
go install github.com/taqnihub/gomigrate@latest
```

</div>

---

## 🚀 Quick Start

```bash
# 1. Install
go install github.com/taqnihub/gomigrate@latest

# 2. Initialize gomigrate in your project
gomigrate init

# 3. Create your first migration
gomigrate create add_users_table

# 4. Apply it
gomigrate up
```

That's it. No Docker commands. No 200-character invocations. Just works.

---

## ✨ Why GoMigrate?

Database migrations don't have to be painful. Compare:

**Before — raw golang-migrate:**
```bash
docker run -it --rm --network host --volume ./db:/db \
  migrate/migrate:v4.19.0 \
  -path=/db/migrations \
  -database "mysql://root:password@localhost:3306/ecomm" \
  up
```

**After — gomigrate:**
```bash
gomigrate up
```

### Features

- 🎨 **Beautiful output** — colored, styled, easy to read
- 🖥️ **Interactive TUI** — run `gomigrate` with no args for a menu
- 📝 **YAML config** — `.gomigrate.yml` keeps credentials out of shell history
- 🔐 **Env var support** — override any setting with `GOMIGRATE_*` variables
- 🛡️ **Safety first** — duplicate detection, dirty-state recovery, clear errors
- 📦 **Library mode** — import `github.com/taqnihub/gomigrate` in your Go code
- 🐬 🐘 **MySQL & PostgreSQL** — full support for both

---

## 📥 Installation

### Install the CLI

With Go 1.22+ installed, one command:

```bash
go install github.com/taqnihub/gomigrate@latest
```

Verify it works:

```bash
gomigrate --version
```

> **Note:** Make sure `$(go env GOPATH)/bin` is in your `$PATH`.
> Add this to your shell config if `gomigrate` isn't found after install:
>
> ```bash
> export PATH=$PATH:$(go env GOPATH)/bin
> ```

### Use as a Go library

Add GoMigrate to your Go project:

```bash
go get github.com/taqnihub/gomigrate@latest
```

Then import in your code:

```go
import (
    "github.com/taqnihub/gomigrate/config"
    "github.com/taqnihub/gomigrate/migrate"
)
```

See [Library Usage](#-library-usage) below for full examples.

### For contributors

Only needed if you want to modify the code or submit a PR:

```bash
git clone https://github.com/taqnihub/gomigrate.git
cd gomigrate
go build -o gomigrate .
./gomigrate --help
```

---

## 📖 CLI Usage

### Interactive mode

Just run `gomigrate` with no arguments:

```bash
gomigrate
```

You'll get a menu to navigate with arrow keys.

### Direct commands

```bash
gomigrate init                         # Set up .gomigrate.yml
gomigrate create add_users_table       # Create new migration files
gomigrate up                           # Apply all pending migrations
gomigrate up 1                         # Apply only next 1 migration
gomigrate down                         # Revert last migration
gomigrate down 3                       # Revert last 3 migrations
gomigrate down --all                   # Revert all migrations
gomigrate status                       # Show migration status
gomigrate force 20260417042058         # Force version (fix dirty state)
```

### Command reference

| Command | Description |
|---------|-------------|
| `init` | Create `.gomigrate.yml` interactively |
| `create <n>` | Create up/down SQL files |
| `up [n]` | Apply all pending (or next N) migrations |
| `down [n]` | Revert last migration (or last N) |
| `status` | Show applied vs pending migrations |
| `force <version>` | Reset version (for dirty states) |

---

## 🔧 Configuration

GoMigrate reads config from (in priority order):

1. CLI flags (`--driver mysql`, `--dir ./migrations`, etc.)
2. Environment variables (`GOMIGRATE_*`)
3. `.gomigrate.yml` in the current directory
4. `.gomigrate.yml` in your home directory
5. Sensible defaults

### Config file

`.gomigrate.yml`:

```yaml
driver: mysql              # mysql or postgres
host: localhost
port: 3306
database: myapp
user: root
password: secret
migrations_dir: ./db/migrations

# Optional
ssl_mode: disable          # disable, require, verify-full
timezone: UTC
lock_timeout: 15s
```

### Environment variables

Perfect for CI/CD:

```bash
export GOMIGRATE_DRIVER=postgres
export GOMIGRATE_HOST=prod-db.example.com
export GOMIGRATE_PASSWORD=$DB_PASSWORD

gomigrate up
```

---

## 📚 Library Usage

GoMigrate is also a Go library. Import it in your code:

```go
package main

import (
    "log"

    "github.com/taqnihub/gomigrate/config"
    "github.com/taqnihub/gomigrate/migrate"
)

func main() {
    // Load config from .gomigrate.yml or env vars
    cfg, err := config.Load("")
    if err != nil {
        log.Fatal(err)
    }

    // Create engine
    engine, err := migrate.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer engine.Close()

    // Apply all pending migrations
    if err := engine.Up(); err != nil {
        log.Fatal(err)
    }

    // Check current version
    version, dirty, _ := engine.Version()
    log.Printf("at version %d (dirty: %v)", version, dirty)
}
```

### Common use case: auto-migrate on app startup

```go
func main() {
    cfg, _ := config.Load("")
    engine, _ := migrate.New(cfg)
    defer engine.Close()

    // Auto-apply pending migrations when the app starts
    if err := engine.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
        log.Fatal("migration failed:", err)
    }

    // ... start your web server ...
}
```

---

## 📝 Writing Migrations

Migration files are plain SQL. GoMigrate creates pairs of files:

```
db/migrations/
├── 20260417100000_add_users_table.up.sql
└── 20260417100000_add_users_table.down.sql
```

**`up.sql`** — apply the change:

```sql
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**`down.sql`** — undo the change:

```sql
DROP TABLE IF EXISTS users;
```

### Naming conventions

| Pattern | Example |
|---------|---------|
| `create_<table>_table` | `create_users_table` |
| `add_<column>_to_<table>` | `add_email_to_users` |
| `drop_<column>_from_<table>` | `drop_legacy_id_from_orders` |
| `add_index_<n>` | `add_index_email_to_users` |

---

## 🆚 Comparison

| Feature | gomigrate | golang-migrate CLI | goose |
|---------|:---------:|:------------------:|:-----:|
| Interactive TUI | ✅ | ❌ | ❌ |
| YAML config file | ✅ | ❌ | ❌ |
| Env var support | ✅ | ❌ | ✅ |
| Pretty colored output | ✅ | ❌ | ⚠️ |
| Duplicate detection | ✅ | ❌ | ❌ |
| Library mode | ✅ | ✅ | ✅ |
| MySQL support | ✅ | ✅ | ✅ |
| PostgreSQL support | ✅ | ✅ | ✅ |

---

## 🛠️ Development (for contributors)

### Prerequisites

- Go 1.22 or higher
- MySQL or PostgreSQL (for testing)

### Clone and build

```bash
git clone https://github.com/taqnihub/gomigrate.git
cd gomigrate
go mod download
make build
```

### Available commands

```bash
make build        # Build the binary
make test         # Run unit tests
make test-cover   # Run tests with coverage report
make install      # Install to $GOPATH/bin
make clean        # Remove build artifacts
```

### Project structure

```
gomigrate/
├── cmd/           # Cobra CLI commands
├── config/        # Config loading (public library package)
├── migrate/       # Migration engine (public library package)
├── internal/tui/  # TUI components (private)
└── main.go        # CLI entry point
```

---

## 🤝 Contributing

Contributions welcome! Please:

1. Fork the repo
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. Make your changes
4. Add tests
5. Run `make test` to verify
6. Commit with conventional messages (`feat:`, `fix:`, `docs:`, etc.)
7. Push and open a PR

### Reporting issues

Found a bug? Have a feature request? [Open an issue](https://github.com/taqnihub/gomigrate/issues).

---

## 📄 License

MIT © [taqnihub](https://github.com/taqnihub)

---

## 🙏 Credits

Built with:

- [golang-migrate](https://github.com/golang-migrate/migrate) — the powerful migration engine under the hood
- [cobra](https://github.com/spf13/cobra) — CLI framework
- [viper](https://github.com/spf13/viper) — config management
- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling
- [huh](https://github.com/charmbracelet/huh) — interactive forms

<div align="center">

**If gomigrate saves you time, consider giving it a ⭐ on GitHub!**

</div>