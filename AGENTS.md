# Repository Guidelines

## Project Structure & Module Organization

This is a **multi-module Go monorepo** providing reusable backend libraries. Each top-level directory is an independent Go module with its own `go.mod`:

```
go-kit/
  dafi/          # Query parsing: filtering, sorting, pagination from URL params
  database/      # PostgreSQL pool (pgx), context-routed transactions, work unit
  di/            # Dependency injection wrapper around samber/do
  httpresponse/  # HTTP response helpers and RFC 9457 problem details
  logger/        # Logger interface with slog adapter and no-op implementation
  server/        # HTTP server with health checks, CORS, recovery middleware
  sqlcraft/      # Fluent SQL query builder (depends on dafi)
```

Inter-module dependencies use `replace` directives pointing to sibling directories (e.g., `../dafi`). There is no workspace `go.work` file.

## Build, Test, and Development Commands

Each module is tested independently from its own directory:

```sh
cd <module>       # e.g. cd sqlcraft
go test ./...     # run all tests in the module
go vet ./...      # static analysis
```

To test everything from the repo root:

```sh
for d in dafi database di httpresponse logger server sqlcraft; do (cd "$d" && go test ./...); done
```

There is no Makefile, linter config, or CI pipeline checked in.

## Coding Style & Naming Conventions

- **Go version**: 1.25.0 across all modules.
- **File naming**: `snake_case.go`. Tests use the `_test.go` suffix.
- **Package layout**: each module has a `core.go` for central types/constructors, plus descriptive files per concern (e.g., `filter.go`, `sort.go`).
- **Exported names**: PascalCase. Interfaces end with descriptive nouns (e.g., `PoolInterface`, `WorkUnit`, `HealthChecker`).
- **Methods**: value-receiver fluent builders (see `dafi.Criteria`), pointer receivers for stateful types.
- **No code generation or linter tooling** is configured. Follow standard `gofmt` formatting.

## Testing Guidelines

- **Framework**: `github.com/stretchr/testify` (`assert` / `require`) in every module.
- **Pattern**: table-driven tests with `t.Run` subtests. See `dafi/parse_test.go` for the canonical example.
- **Test naming**: `Test<Type>_<Method>` or `Test<Function>` (e.g., `TestQueryParser_Parse`).
- **Mocking**: interfaces are defined for external dependencies (`PoolInterface`, `DatabasePort`) to enable test doubles.

## Commit & Pull Request Guidelines

Commit messages use **lowercase imperative style**, prefixed with an action verb and the affected package scope:

```
add <package> package with <description>
refactor <package>/<file>: <what changed>
```

Keep commits focused on a single package or concern. Reference related packages when cross-cutting changes are needed.
