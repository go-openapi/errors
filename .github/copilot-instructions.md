# Copilot Instructions for go-openapi/errors

## Project Overview

This package provides structured error types for the go-openapi/go-swagger ecosystem.
It defines an `Error` interface (with an HTTP status code and message) and concrete types
for validation failures, parsing errors, authentication errors, and HTTP middleware errors.

### Package layout

| File | Contents |
|------|----------|
| `api.go` | Core `Error` interface, constructors (`New`, `NotFound`, `NotImplemented`, `MethodNotAllowed`), `ServeError` HTTP handler |
| `schema.go` | `Validation` struct, `CompositeError`, ~30 validation constructors, error code constants |
| `headers.go` | `InvalidContentType`, `InvalidResponseFormat` |
| `parsing.go` | `ParseError` struct for parameter parsing failures |
| `auth.go` | `Unauthenticated` constructor |
| `middleware.go` | `APIVerificationFailed` for spec/registration mismatches |

### Key design points

- Error codes >= 600 are domain codes (not HTTP); `ServeError` maps them to 422.
- All error types implement `json.Marshaler` for structured JSON responses.
- Zero runtime dependencies outside the Go standard library.

## Conventions

Coding conventions are found beneath `.github/copilot`

### Summary

- All `.go` files must have SPDX license headers (Apache-2.0).
- Commits require DCO sign-off (`git commit -s`).
- Linting: `golangci-lint run` — config in `.golangci.yml` (posture: `default: all` with explicit disables).
- Every `//nolint` directive **must** have an inline comment explaining why.
- Tests: `go test ./...`. CI runs on `{ubuntu, macos, windows} x {stable, oldstable}` with `-race`.
- Test framework: `github.com/go-openapi/testify/v2` (not `stretchr/testify`; `testifylint` does not work).

See `.github/copilot/` (symlinked to `.claude/rules/`) for detailed rules on Go conventions, linting, testing, and contributions.
