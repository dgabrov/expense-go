# Repository Guidelines

## Project Structure & Module Organization

This repository is a Go HTTP service. `main.go` starts the application. Application setup, routing, and database connection code live in `internal/proj`; HTTP handlers are in `internal/proj/controller`, and SQL operations belong on `internal/proj/server.Server`. Shared request and domain types are in `internal/data`. Database schema and migration material are in `docs/db`, API specifications are in `docs/spec`, and `docs/info` contains configuration examples. Keep implementation packages under `internal/`.

## Build, Test, and Development Commands

- `go build ./...` compiles all packages.
- `go test ./...` runs the Go test suite.
- `go vet ./...` checks for common Go mistakes.
- `go run .` starts the service. Set `CONFIG_FILE` to the path of a JSON configuration file before running.

Format Go changes with `gofmt` (for example, `gofmt -w internal/proj`). No separate formatter or linter configuration is present.

## Coding Style & Naming Conventions

Follow standard Go formatting and naming: use `MixedCaps` for exported identifiers, concise lower-case names for packages, and descriptive file names such as `post_add_transaction.go`. Keep HTTP handling in controllers and database queries or mutations in methods on `server.Server`. Use context-aware database calls and parameterized SQL, consistent with existing code.

## Testing Guidelines

Use Go's built-in `testing` package and place tests beside the code they cover in `*_test.go` files. Run `go test ./...` after changes. There are currently no test files, so add focused tests when introducing behavior, especially for handlers and database operations.

## Commit & Pull Request Guidelines

Recent commits use short action-oriented subjects, including `Add POST /addBulk endpoint with all-or-none transaction insert` and `Set utf8mb4 charset and collation on all DDL tables`. Use a concise imperative subject describing the change. Pull requests should summarize behavior and affected areas, link related issues when available, and report relevant build, vet, and test results. Include request/response examples when changing an API endpoint.

## Configuration & Database

Keep credentials in the external JSON file selected by `CONFIG_FILE`; do not commit secrets. The database schema of record is `docs/db/db_ddl.sql`. Update schema or migration documentation when a database change requires it.
