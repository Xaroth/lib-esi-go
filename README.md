# lib-esi-go

`lib-esi-go` provides generated ESI request packages on top of standard `net/http` interfaces.

Build an `http.Client` with an ESI-aware `http.RoundTripper`, then execute generated requests from `esi/...`.

## Simple Usage

In normal use, `transport.New(...)` becomes the client's `Transport`, and the generated request package supplies the typed `Request(...)` function.

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/esi/getmetastatus"
	"github.com/xaroth/lib-esi-go/transport"
)

func main() {
	client := &http.Client{
		Transport: transport.New(
			"my-app",
			"1.0.0",
			[]string{"mailto:esi@example.com"},
			defaults.CompatibilityDate,
		),
	}

	resp, err := getmetastatus.Request(context.Background(), client)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= 400 {
		panic(resp.Status)
	}

	fmt.Printf("routes: %d\n", len(resp.Data.Routes))
}
```

Generated packages under `esi/...` are thin typed bindings. They accept any sender that implements:

```go
type RequestSender interface {
	Do(req *http.Request) (*http.Response, error)
}
```

In practice, that is usually `*http.Client`.

## Authenticated Requests

Authentication is request-scoped. The transport always includes the authentication middleware, but it only adds an `Authorization` header when the request carries a token.

### Token Interface

Provide any type that implements `authentication.Token`:

```go
type Token interface {
	Owner() int64
	Token() string
}
```

`Owner()` identifies the token owner. This matters for middleware such as rate limiting, which can bucket requests per owner.

Pass the token to a request with `authentication.WithToken(...)`.

### Just-In-Time Refresh

If your token type also implements `authentication.RefreshableToken`, the authentication middleware calls `RefreshIfNeeded(ctx)` immediately before it sets the `Authorization` header.

That keeps the token valid even if the request is delayed by middleware, queued behind rate limiting, or re-issued by your own retry logic.

### Example

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/common/character"
	"github.com/xaroth/lib-esi-go/esi/getcharacterscharacteridlocation"
	"github.com/xaroth/lib-esi-go/middleware/authentication"
	"github.com/xaroth/lib-esi-go/transport"
)

type staticToken struct {
	owner int64
	value string
}

func (t staticToken) Owner() int64  { return t.owner }
func (t staticToken) Token() string { return t.value }

func main() {
	client := &http.Client{
		Transport: transport.New(
			"my-app",
			"1.0.0",
			[]string{"mailto:esi@example.com"},
			defaults.CompatibilityDate,
		),
	}

	resp, err := getcharacterscharacteridlocation.Request(
		context.Background(),
		client,
		&getcharacterscharacteridlocation.Input{
			Character: character.Identifier(123456789),
		},
		authentication.WithToken(staticToken{
			owner: 123456789,
			value: "<access-token>",
		}),
	)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode >= 400 {
		panic(resp.Status)
	}

	fmt.Printf("solar system: %d\n", resp.Data.SolarSystemId)
}
```

## Middlewares

Middleware lives at the `http.RoundTripper` layer. `transport.New(...)` builds a transport chain with the default ESI middleware, and `transport.WithMiddleware(...)` appends additional middleware.

Use middleware for cross-cutting transport behavior such as headers, auth, caching, rate limiting, and logging. Keep request-specific input and output logic in the generated `esi/...` packages.

### Built In By Default

`transport.New(...)` always includes:

- compatibility date
- language
- timeout
- user agent
- tenant
- authentication

The following snippets show the request call only. Imports and surrounding client setup are omitted.

#### Compatibility Date

Configured as the `compatibilityDate` argument to `transport.New(...)`.

```go
rt := transport.New("my-app", "1.0.0", contacts, "2026-05-19")
```

Override or clear it per request with `compatibilitydate.WithCompatibilityDate(...)`.

```go
resp, err := getmetastatus.Request(
	ctx,
	client,
	compatibilitydate.WithCompatibilityDate("2026-05-19"),
)
```

Passing an empty string omits the header for that request.

#### Language

By default, `transport.New(...)` uses `defaults.Language`.

Override or clear it per request with `language.WithLanguage(...)`.

```go
resp, err := getmetastatus.Request(
	ctx,
	client,
	language.WithLanguage("en"),
)
```

#### Timeout

By default, `transport.New(...)` uses `defaults.RequestTimeout`.

Override or disable it per request with `timeout.WithTimeout(...)`.

```go
resp, err := getmetastatus.Request(
	ctx,
	client,
	timeout.WithTimeout(5*time.Second),
)
```

A timeout value less than or equal to zero disables the timeout for that request.

#### User-Agent

Configured by the first three arguments to `transport.New(...)`:

- application name
- application version
- contact values

```go
rt := transport.New(
	"my-app",
	"1.0.0",
	[]string{"mailto:esi@example.com", "https://example.com"},
	defaults.CompatibilityDate,
)
```

#### Tenant

By default, `transport.New(...)` uses `defaults.Tenant`.

Override or clear it per request with `tenant.WithTenant(...)`.

```go
resp, err := getmetastatus.Request(
	ctx,
	client,
	tenant.WithTenant("tranquility"),
)
```

### Not Enabled By Default

#### Cache

The cache middleware is opt-in:

```go
package main

import (
	"net/http"

	_ "github.com/glebarez/go-sqlite"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/middleware/cache"
	"github.com/xaroth/lib-esi-go/transport"
)

func newCachedClient() *http.Client {
	return &http.Client{
		Transport: transport.New(
			"my-app",
			"1.0.0",
			[]string{"mailto:esi@example.com"},
			defaults.CompatibilityDate,
			transport.WithMiddleware(cache.Middleware("./cache.sqlite")),
		),
	}
}
```

The cache storage uses SQLite by default, so you must include the side-effect import for `github.com/glebarez/go-sqlite`.

#### Rate Limiting

The rate-limiting middleware is also opt-in:

```go
package main

import (
	"net/http"

	defaults "github.com/xaroth/lib-esi-go"
	"github.com/xaroth/lib-esi-go/middleware/ratelimiting"
	"github.com/xaroth/lib-esi-go/middleware/ratelimiting/memory"
	"github.com/xaroth/lib-esi-go/transport"
)

func newRateLimitedClient() *http.Client {
	return &http.Client{
		Transport: transport.New(
			"my-app",
			"1.0.0",
			[]string{"mailto:esi@example.com"},
			defaults.CompatibilityDate,
			transport.WithMiddleware(ratelimiting.Middleware(memory.New())),
		),
	}
}
```

The built-in backend is `middleware/ratelimiting/memory`. It tracks ESI rate limit headers and delays requests when usage approaches the configured target. If you prefer using your own (distributed) rate limiting logic, implement `ratelimiting.RateLimiter` yourself.

### Custom Middleware

Custom middleware implements `middleware.Middleware`:

```go
type Middleware func(http.RoundTripper) http.RoundTripper
```

The helper type `middleware.MiddlewareFunc` makes it easy to wrap a `RoundTripper`:

```go
package main

import (
	"log"
	"net/http"

	"github.com/xaroth/lib-esi-go/middleware"
)

func loggingMiddleware(next http.RoundTripper) http.RoundTripper {
	return middleware.MiddlewareFunc(func(req *http.Request) (*http.Response, error) {
		log.Printf("%s %s", req.Method, req.URL.Path)
		return next.RoundTrip(req)
	})
}
```

Add it with `transport.WithMiddleware(...)`, or use `transport.NewChain(...)` if you want full control over the chain.

## Code Generation

Common models and ESI requests are generated in this repository:

The generated outputs in this repository live under:

- `common/...`
- `esi/...`

These folders are periodically updated (as well as the compatibility date in defaults.go) based on the latest
available ESI compatibility date.

If you wish to generate your own common models and/or requests, have a look at the `cmd` directory.
