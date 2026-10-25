# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

- Build: `go build ./...`
- Vet: `go vet ./...`
- Test all: `go test ./...`
- Single test: `go test ./internal/<pkg> -run '^TestName$'`
- Run: `go run .`

The module is `github.com/dgabrov/expense-go` (Go 1.26, toolchain 1.26.2 is installed locally). There is no Makefile, CI config, linter config, or README yet.

## Architecture

This is an early scaffold with no real domain logic yet.

- `main.go` is the entry point. It calls `proj.Start()` and logs any error with `slog`.
- `internal/proj` holds application startup. `Start()` loads the externalized config first. `config.go` reads the file named by the `CONFIG_FILE` env var (relative to the working directory) and decodes it into `data.ConfigData`.
- `internal/data` defines the data types. It currently has only `ConfigData{Version string}`.
- `docs/info/config.json` holds a version record (`"version": "1.0"`) that is not yet read by code. `docs/db/` and `docs/spec/` are empty placeholders.

The `internal/` packages are where new application code should go. `internal/data` is the natural home for shared types, and `internal/proj` is the natural home for wiring them together at startup.

## General

- Any new acces to db, you add new method to server.Server
- never push git nor ask about it
- the database is in @docs/db/db_ddl.sql and you can connect locally at localhost:3306 login, password database all "expense", mariaDB
