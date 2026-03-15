# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This package provides structured error types for the go-openapi/go-swagger ecosystem.
It defines an `Error` interface (with an HTTP status code and message) and concrete types
for validation failures, parsing errors, authentication errors, and HTTP middleware errors.
These error types are consumed by validators, spec loaders, and API runtimes to report
problems back to API consumers in a structured, JSON-serializable way.

See [docs/MAINTAINERS.md](../docs/MAINTAINERS.md) for CI/CD, release process, and repo structure details.

### Package layout

| File | Contents |
|------|----------|
| `doc.go` | Package-level godoc |
| `api.go` | Core `Error` interface, `apiError` struct, `New`/`NotFound`/`NotImplemented`/`MethodNotAllowed` constructors, `ServeError` HTTP handler, `CompositeError` flattening |
| `schema.go` | `Validation` struct, `CompositeError` struct, and ~30 constructor functions for JSON Schema / OpenAPI validation failures (type, required, enum, min/max, pattern, etc.); error code constants (`InvalidTypeCode`, `RequiredFailCode`, ...) |
| `headers.go` | `Validation` constructors for HTTP header errors: `InvalidContentType`, `InvalidResponseFormat` |
| `parsing.go` | `ParseError` struct and `NewParseError` constructor for parameter parsing failures |
| `auth.go` | `Unauthenticated` constructor (401) |
| `middleware.go` | `APIVerificationFailed` struct for mismatches between spec and registered handlers |

### Key API

- `Error` interface -- `error` + `Code() int32`
- `New(code, message, args...)` -- general-purpose error constructor
- `NotFound` / `NotImplemented` / `MethodNotAllowed` / `Unauthenticated` -- HTTP error constructors
- `Validation` struct -- carries `Name`, `In`, `Value` context for validation failures
- ~30 validation constructors: `Required`, `InvalidType`, `TooLong`, `TooShort`, `EnumFail`, `ExceedsMaximum`, `FailedPattern`, `DuplicateItems`, `TooManyItems`, `TooFewItems`, `PropertyNotAllowed`, `CompositeValidationError`, etc.
- `CompositeError` -- groups multiple errors; implements `Unwrap() []error`
- `ServeError(rw, r, err)` -- HTTP handler that serializes any `Error` to JSON
- `ParseError` -- parsing failure with `Name`, `In`, `Value`, `Reason`

### Dependencies

- `github.com/go-openapi/testify/v2` -- test-only assertions (zero-dep testify fork)

This package has **zero runtime dependencies** outside the Go standard library.

### Notable design decisions

- **Error codes above 599 are domain codes, not HTTP codes** -- validation error codes (`InvalidTypeCode = 600`, `RequiredFailCode = 601`, ...) use values >= 600 so they can be distinguished from real HTTP status codes. `ServeError` maps any code >= 600 to `DefaultHTTPCode` (422).
- **`Validation.ValidateName` mutates in place** -- callers chain `.ValidateName(name)` to prepend parent property names for nested validation errors, building dotted paths like `address.street`.
- **All error types implement `json.Marshaler`** -- errors serialize to structured JSON with `code`, `message`, and type-specific fields (`name`, `in`, `value`), not just a string.
- **`CompositeError` flattens recursively** -- `ServeError` flattens nested `CompositeError` trees and serves only the first leaf error to avoid overwhelming API consumers.
