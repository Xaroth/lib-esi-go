# EVE Online ESI Go Library - Implementation Plan

## Overview

This project creates a Go library for interacting with EVE Online's ESI (EVE Swagger Interface) API. The library uses a code generation approach where a CLI tool (`esi-gen`) generates type-safe Go packages from the OpenAPI specification.

## Architecture

### Single Approach: Code Generation Tool

- **`esi-gen` CLI Tool**: Generates all packages from OpenAPI spec
- **Generated Packages**: Type-safe Go packages for ESI API endpoints
- **Core Library Package**: Common interfaces and utilities
- **No runtime dependencies** on the OpenAPI spec

## Project Structure

```
lib-esi-go/
├── cmd/
│   └── esi-gen/              # Code generation CLI tool
│       ├── main.go
│       ├── generator/        # Code generation logic
│       ├── templates/        # Go templates for generation
│       └── fetcher/          # OpenAPI spec fetching
├── client/               # HTTP client & middleware
├── ratelimit/            # Floating window rate limiter
├── auth/                 # OAuth2 interfaces
├── cache/                # Generic caching interfaces
│── go.mod
├── go.sum
└── README.md
```

## Technology Stack

- **HTTP Client**: `net/http` (standard library)
- **Logging**: `log/slog` (Go 1.21+ structured logging)
- **OAuth2**: `golang.org/x/oauth2`
- **Code Generation**: `oapi-codegen` + custom templates
- **CLI**: `flag` package (standard library)
- **Testing**: `testing` + `testify/assert` for table-driven tests

## Key Features

1. **Zero Runtime Dependencies** on OpenAPI spec
2. **Generated Packages** for type-safe API access
3. **Floating Window Rate Limiting** with ESI-specific headers
4. **Generic Interfaces** for auth, caching, logging
5. **Go 1.25 JSON v2** optional support via build tags
6. **Compatibility Date Management** first-class support (header or query param)
7. **User Agent** formatting for ESI ecosystem citizenship

## Implementation Phases

### Phase 1: Core Library Foundation

#### 1.1 Project Setup
- Initialize Go module with Go 1.25
- Set up basic project structure
- Add core dependencies

#### 1.2 Core Interfaces
```go
// Authentication interface
type Authenticator interface {
    GetToken(ctx context.Context) (string, error)
    RefreshToken(ctx context.Context) (string, error)
}

// Rate limiter interface
type RateLimiter interface {
    Wait(ctx context.Context, group string) error
    UpdateLimits(headers http.Header)
}

// Cache interface
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, data []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

// Logger interface
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
}

// User agent configuration
type UserAgentConfig struct {
    LibraryName    string
    LibraryVersion string
    ContactInfo    string // email or GitHub username
}
```

#### 1.3 HTTP Client with Middleware
- Base HTTP client using `net/http`
- Middleware system for auth, rate limiting, logging
- Request/response interceptors
- Error handling; pluggable retry hooks (no automatic retries)
- Proper User-Agent header formatting

### Phase 2: Rate Limiting Implementation

#### 2.1 ESI Rate Limiter
ESI uses floating window rate limiting with:
- Potentially multiple rate limit groups (surfaced via OpenAPI `x-rate-limit` extension when present)
- Token bucket system with configurable window sizes
- Headers: `X-Rate-Limit-Remaining`, `X-Rate-Limit-Reset`

```go
type ESIRateLimiter struct {
    groups map[string]*TokenBucket
    mu     sync.RWMutex
}

type TokenBucket struct {
    maxTokens   int
    tokens      int
    windowSize  time.Duration
    lastRefill  time.Time
}
```

#### 2.2 Rate Limit Middleware
- Automatic rate limiting based on endpoint groups
- Header parsing and limit updates
- Context-aware waiting
- Consume `x-rate-limit` OpenAPI extensions when available; fall back to conservative defaults

### Phase 3: Code Generation Tool

#### 3.1 `esi-gen` CLI Tool
```bash

# Generate all packages
esi-gen generate --compatibility-date YYYY-MM-DD --tier tranquility --output ./generated/ --package generated

# Generate specific tags
esi-gen generate --compatibility-date YYYY-MM-DD --tier tranquility --output ./generated/ --package generated tag1 tag2 tag3
```

#### 3.2 Code Generation Features
- Fetch OpenAPI spec from ESI endpoint
- Support compatibility date selection (header or `compatibility_date` query parameter)
- Generate packages using `oapi-codegen` with custom templates
- Include rate limiting metadata from `x-rate-limit` extensions where present
- Optional Go 1.25 JSON v2 support via build tags
- Support ESI tier selection for spec fetching (e.g., `tranquility`, `singularity`)

#### 3.3 Generated Package Structure
Each generated package will contain:
- Type-safe request/response models
- Client methods for each endpoint
- Rate limiting group information
- Proper error handling
- Consistent User-Agent header formatting

### Phase 4: Go 1.25 JSON v2 Integration

#### 4.1 JSON v2 Configuration
```go
//go:build go1.25

package yourpkg

import (
    json "encoding/json/v2"
)
```

#### 4.2 Build Tags
- Use build tags to support both JSON v1 and v2
- Graceful fallback for older Go versions
- Performance improvements with v2
```go
//go:build !go1.25

package yourpkg

import (
    json "encoding/json"
)
```

### Phase 5: Testing Strategy

#### 5.1 Unit Tests
- Table-driven tests for all components
- Mock implementations for interfaces
- Rate limiting logic tests
- Code generation tests

#### 5.2 Integration Tests
- CI/CD pipeline with ESI API tests
- Compatibility date validation
- Rate limiting behavior verification
- Authentication flow tests

#### 5.3 Test Structure
```go
func TestRateLimiter(t *testing.T) {
    tests := []struct {
        name     string
        group    string
        headers  http.Header
        expected time.Duration
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## ESI-Specific Considerations

### Rate Limiting
- ESI uses floating window rate limiting
- Different limits for different endpoint groups
- Headers: `X-Rate-Limit-Remaining`, `X-Rate-Limit-Reset`
- `x-rate-limit` extensions in the OpenAPI spec denote groups/limits where available

### Authentication
- OAuth2 required for most endpoints
- Token refresh handling
- Scope-based access control

### Versioning
- Compatibility date via `X-Compatibility-Date` header (or `compatibility_date` query parameter)
- ESI rejects future dates; date rolls over at 11:00 UTC
- Response includes `X-Compatibility-Date` with the matched date
- Prefer pinning generator inputs to a specific date for reproducibility

### Caching
- ETag support with `If-None-Match` headers
- Cache-Control headers
- Generic caching interface; pass through conditional request headers provided by users and surface response cache metadata
- Optional in-memory cache implementation for convenience

### User Agent
- Proper User-Agent header formatting for ESI ecosystem citizenship
- Format: `<ServiceName>/<ServiceVersion> (<contact-info>) esi-lib-go/<lib-version> (<lib-contact-info>)
- Contact info should include at least one of: github user, email address, or EVE Character name
- Library info is hard-coded, and cannot be changed
- Helps ESI developers with diagnostics and rate limiting
- Example: `third-party-application/1.2.3 (github.com/app-creator, eve:App Creator) esi-lib-go/1.0.0 (github.com/Xaroth)`

## Development Workflow

### 1. Initial Setup
```bash
# Clone and setup
git clone <repo>
cd lib-esi-go

# Install dependencies
go mod tidy

# Build esi-gen tool
go build ./cmd/esi-gen
```

### 2. Code Generation
```bash
# Fetch latest spec
./esi-gen fetch --output ./specs/

# Generate all packages
./esi-gen generate --spec ./specs/openapi.yaml --output ./generated/

# Run tests
go test ./...
```

### 3. Development
- Modify core library
- Update templates in `cmd/esi-gen/templates/`
- Regenerate packages as needed
- Run tests continuously

## Dependencies

### Core Dependencies
- `golang.org/x/oauth2` - OAuth2 support
- `github.com/getkin/kin-openapi` - OpenAPI spec parsing
- `github.com/deepmap/oapi-codegen` - Code generation

### Development Dependencies
- `github.com/stretchr/testify` - Testing utilities
- `github.com/stretchr/testify/mock` - Mock generation

### Build Dependencies
- `golang.org/x/tools` - Go tooling support

## Future Enhancements

1. **Caching Backends** - Redis, in-memory, file-based
2. **Metrics** - Prometheus metrics for rate limiting
3. **Response Validation** - Schema validation for responses

## Success Criteria

1. **Type Safety** - All API calls are type-safe
2. **Performance** - Efficient rate limiting and HTTP handling
3. **Usability** - Simple, intuitive API for library users
4. **Maintainability** - Easy to update when ESI changes
5. **Testing** - Comprehensive test coverage
6. **Documentation** - Clear examples and API documentation

## Timeline

- **Week 1-2**: Core library foundation and interfaces
- **Week 3-4**: Rate limiting implementation and HTTP client
- **Week 5-6**: Code generation tool and templates
- **Week 7-8**: Package generation and testing
- **Week 9-10**: Integration testing and documentation
- **Week 11-12**: Polish, optimization, and release preparation
