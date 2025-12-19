# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a shared errors library for the go-openapi toolkit. It provides an `Error` interface and concrete error types for API errors and JSON-schema validation errors. The package is used throughout the go-openapi ecosystem (github.com/go-openapi).

## Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Run a specific test
go test -v -run TestName ./...
```

### Linting
```bash
# Run golangci-lint (must be run before committing)
golangci-lint run
```

### Building
```bash
# Build the package
go build ./...

# Verify dependencies
go mod verify
go mod tidy
```

## Architecture and Code Structure

### Core Error Types

The package provides a hierarchy of error types:

1. **Error interface** (api.go:20-24): Base interface with `Code() int32` method that all errors implement
2. **apiError** (api.go:26-37): Simple error with code and message
3. **CompositeError** (schema.go:94-122): Groups multiple errors together, implements `Unwrap() []error`
4. **Validation** (headers.go:12-55): Represents validation failures with Name, In, Value fields
5. **ParseError** (parsing.go:12-42): Represents parsing errors with Reason field
6. **MethodNotAllowedError** (api.go:74-88): Special error for method not allowed with Allowed methods list
7. **APIVerificationFailed** (middleware.go:12-39): Error for API spec/registration mismatches

### Error Categorization by File

- **api.go**: Core error interface, basic error types, HTTP error serving
- **schema.go**: Validation errors (type, length, pattern, enum, min/max, uniqueness, properties)
- **headers.go**: Header validation errors (content-type, accept)
- **parsing.go**: Parameter parsing errors
- **auth.go**: Authentication errors
- **middleware.go**: API verification errors

### Key Design Patterns

1. **Error Codes**: Custom error codes >= 600 (maximumValidHTTPCode) to differentiate validation types without conflicting with HTTP status codes
2. **Conditional Messages**: Most constructors have "NoIn" variants for errors without an "In" field (e.g., tooLongMessage vs tooLongMessageNoIn)
3. **ServeError Function** (api.go:147-201): Central HTTP error handler using type assertions to handle different error types
4. **Flattening**: CompositeError flattens nested composite errors recursively (api.go:108-134)
5. **Name Validation**: Errors can have their Name field updated for nested properties via ValidateName methods

### JSON Serialization

All error types implement `MarshalJSON()` to provide structured JSON responses with code, message, and type-specific fields.

## Testing Practices

- Uses forked `github.com/go-openapi/testify/v2` for minimal test dependencies
- Tests follow pattern: `*_test.go` files next to implementation
- Test files cover: api_test.go, schema_test.go, middleware_test.go, parsing_test.go, auth_test.go

## Code Quality Standards

### Linting Configuration
- Enable all golangci-lint linters by default, with specific exclusions in .golangci.yml
- Complexity threshold: max 20 (cyclop, gocyclo)
- Line length: max 180 characters
- Run `golangci-lint run` before committing

### Disabled Linters (and why)
Key exclusions from STYLE.md rationale:
- depguard: No import constraints enforced
- funlen: Function length not enforced (cognitive complexity preferred)
- godox: TODOs are acceptable
- nonamedreturns: Named returns are acceptable
- varnamelen: Short variable names allowed when appropriate

## Release Process

- Push semver tag (v{major}.{minor}.{patch}) to master branch
- CI automatically generates release with git-cliff
- Tags should be PGP-signed
- Tag message prepends release notes

## Important Constants

- `DefaultHTTPCode = 422` (http.StatusUnprocessableEntity)
- `maximumValidHTTPCode = 600`
- Custom error codes start at 600+ (InvalidTypeCode, RequiredFailCode, etc.)
