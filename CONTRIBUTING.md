# Contributing to NamedRouter Chi Adapter

Thank you for your interest in contributing to NamedRouter Chi Adapter.

This project is an adapter between
[NamedRouter](https://github.com/mafalt/namedrouter) and
[Chi](https://github.com/go-chi/chi).

The most important principle when contributing is to keep the adapter focused
on translating NamedRouter operations into Chi operations.

## Architecture

The relationship between the projects is:

```text
Application
     │
     ▼
NamedRouter
     │
     ▼
ChiAdapter
     │
     ▼
chi.Router
```

`NamedRouter` defines the common routing API.

`ChiAdapter` implements that API using Chi.

The adapter must not move NamedRouter-specific responsibilities into Chi.

## Responsibilities

### NamedRouter

`NamedRouter` is responsible for:

* route names,
* route registry,
* duplicate route detection,
* reverse routing,
* route parameter validation,
* common routing API.

### ChiAdapter

`ChiAdapter` is responsible for:

* creating and owning the Chi router,
* translating route registration,
* creating Chi subrouters,
* applying Chi middleware,
* joining paths,
* serving static files,
* walking Chi routes,
* parsing Chi route parameters,
* applying Chi route parameters.

If a change requires knowledge of how Chi works internally, it most likely
belongs in this repository.

## Creating the Adapter

The adapter creates its own Chi router:

```go
func New() namedrouter.Adapter {
	return &chiAdapter{
		router:  chi.NewRouter(),
		parser:  &chiParameterParser{},
		applier: &chiParameterApplier{},
	}
}
```

Consumers should not need to create a Chi router themselves.

This is intentional.

The adapter provides the boundary between the application and Chi.

## Parameter Parser

The Chi adapter implements:

```go
type chiParameterParser struct{}
```

which satisfies:

```go
type ParameterParser interface {
	Parse(route RouteDefinition) (RouteParameterNames, error)
}
```

The parser must understand Chi's route parameter syntax.

For example:

```
/users/{id}
/users/{userId}/posts/{postId}
```

The parser extracts the parameter names:

```go
RouteParameterNames{
	"userId": {},
	"postId": {},
}
```

## Keep Parsing Chi-Specific

Do not move Chi parameter syntax into the `NamedRouter` package.

The parser belongs to this adapter because `{parameter}` is part of the
routing syntax understood by Chi.

If Chi changes its parameter syntax, the adapter should be responsible for
adapting to that change.

## Parameter Applier

The Chi adapter also implements:

```go
type chiParameterApplier struct{}
```

which satisfies:

```go
type ParameterApplier interface {
	Apply(pattern string, params RouteParams) string
}
```

The implementation replaces Chi route parameters and URL-escapes their
values.

For example:

```
/users/{id}
```

with:

```go
RouteParams{
	"id": 42,
}
```

becomes:

```
/users/42
```

Path values must be escaped using URL path escaping.

## Subrouters

Chi subrouters are created using Chi's router grouping functionality.

A subrouter must preserve:

* the current route context,
* the supplied prefix,
* subrouter middleware.

Example:

```go
childRouter := c.router.Group(func(r chi.Router) {
	r.Use(middlewares...)
})
```

Nested subrouters must continue to work correctly.

## Middleware

The adapter supports three middleware contexts.

### Global middleware

`Use()` modifies the global Chi router:

```go
c.router.Use(middlewares...)
```

### Subrouter middleware

`Subrouter()` creates a Chi routing context and applies the supplied
middleware to that context.

### Route middleware

`ApplyMiddlewares()` creates the route-specific Chi middleware context.

This operation must affect only the current route registration context.

It must not unexpectedly modify global middleware configuration.

## Route Registration

The adapter receives:

```go
Register(method, pattern string, handler http.Handler)
```

The route name is deliberately not part of the adapter API.

For example:

```go
c.router.Method(method, pattern, handler)
```

The adapter must not maintain its own route-name registry.

Route names belong to `NamedRouter`.

## Path Joining

The adapter implements:

```go
JoinPath(prefix, pattern string) string
```

The implementation currently uses:

```go
path.Join(prefix, pattern)
```

This method is part of the adapter because path semantics are ultimately
part of the underlying router integration.

Do not move Chi-specific path handling into `NamedRouter`.

## Static Files

The adapter implements:

```go
Static(pattern, root string)
```

using Chi together with the standard library's file server functionality.

For example:

```go
c.router.Handle(
	pattern,
	http.StripPrefix(
		stripPrefix,
		http.FileServer(http.Dir(root)),
	),
)
```

Changes to static file handling should preserve the behavior expected from
Chi and `net/http`.

## Route Walking

`Walk()` delegates route inspection to:

```go
chi.Walk(...)
```

The adapter currently prints:

```
METHOD PATH
```

to standard output.

Changes to this behavior should be considered carefully because the
`Walk()` method is part of the common NamedRouter API.

## Testing

Tests should focus on observable behavior.

Prefer testing requests against the adapter rather than inspecting internal
Chi router state.

For example:

```go
req := httptest.NewRequest(
	http.MethodGet,
	"/users/42",
	nil,
)

rec := httptest.NewRecorder()

router.ServeHTTP(rec, req)
```

Tests should cover, where applicable:

* route registration,
* HTTP methods,
* route parameters,
* reverse routing,
* nested subrouters,
* prefixes,
* global middleware,
* subrouter middleware,
* route middleware,
* static files,
* path joining,
* route walking.

### Testing Router-Specific Behavior

The purpose of this repository is to verify that `NamedRouter` semantics are
correctly translated to Chi.

Tests should therefore verify both:

* NamedRouter behavior exposed through the adapter.
* Correct integration with Chi.

Avoid tests that depend unnecessarily on private implementation details.

## Compatibility

This adapter targets Chi v5.

Changes to the Chi dependency should be evaluated carefully.

In particular, consider:

* routing behavior,
* parameter syntax,
* middleware semantics,
* subrouter behavior,
* static file handling.

If a Chi upgrade changes behavior relied upon by the adapter, add regression
tests before changing the implementation.

## Public API Changes

The adapter should expose as little API as possible.

The primary public entry point is:

```go
func New() namedrouter.Adapter
```

Avoid exposing the underlying `chi.Router`.

For example, do not add:

```go
func (c *chiAdapter) Router() chi.Router
```

unless there is a compelling architectural reason.

The purpose of the adapter is to prevent consumers from depending directly
on Chi.

## Dependency Direction

The dependency direction must remain:

```
NamedRouter
    ▲
    │
ChiAdapter
    │
    ▼
  Chi
```

`ChiAdapter` depends on `NamedRouter`.

`NamedRouter` must never depend on `ChiAdapter`.

## Pull Requests

Before submitting a pull request:

* Run the complete test suite.
* Run `go vet ./...`.
* Run `gofmt`.
* Add or update tests for changed behavior.
* Update documentation when the public API or behavior changes.
* Verify that no unnecessary Chi implementation details leak through the adapter API.

## Architecture First

When making a change, ask:

> Is this a `NamedRouter` concern or a Chi concern?

If the answer is Chi-specific, it belongs here.

If the answer is related to named routes, route registry management,
reverse-routing validation, or the common API, it belongs in NamedRouter.

The goal of this project is to keep the adapter thin while preserving the
semantics of the underlying Chi router.

## License

By contributing to this project, you agree that your contributions will be
licensed under the same license as the project.
