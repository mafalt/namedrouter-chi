# NamedRouter Chi Adapter

[![CI](https://github.com/mafalt/namedrouter-chi/actions/workflows/ci.yml/badge.svg)](https://github.com/mafalt/namedrouter-chi/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/mafalt/namedrouter-chi.svg)](https://pkg.go.dev/github.com/mafalt/namedrouter-chi)
[![License](https://img.shields.io/github/license/mafalt/namedrouter-chi)](LICENSE)

Chi adapter for [NamedRouter](https://github.com/mafalt/namedrouter).

This package allows NamedRouter to use [Chi](https://github.com/go-chi/chi)
as its underlying HTTP router.

## Installation

```bash
go get github.com/mafalt/namedrouter-chi
```

## Usage

Create the adapter and pass it to NamedRouter:

```go
package main


import (
	"net/http"


	"github.com/mafalt/namedrouter"
	chiadapter "github.com/mafalt/namedrouter-chi"
)


func main() {
	router := namedrouter.New(chiadapter.New())


	router.RegisterGet(
		"/users/{id}",
		"user",
		getUser,
	)


	http.ListenAndServe(":8080", router)
}


func getUser(w http.ResponseWriter, r *http.Request) {
	// ...
}
```

The adapter creates and owns the underlying Chi router. The application does
not need to create or access a `chi.Router` directly.

## Chi-specific Behavior

Route parameters

The adapter uses Chi's standard route parameter syntax:

```
/users/{id}
/users/{userId}/posts/{postId}
```

These parameters are also supported by NamedRouter's reverse routing:

```go
url := router.MustURL(
	"user",
	namedrouter.RouteParams{
		"id": 42,
	},
)
```

Parameter values are URL path escaped before being inserted into the
generated URL.

`NamedRouter` provides framework-independent access to URL parameters through URLParam:

```go
id := router.URLParam(r, "id")
```

The actual URL parameter lookup is delegated to the configured router adapter, so application code does not need to depend on the underlying routing framework.

For example, the Chi adapter uses Chi's URL parameter mechanism internally:

```go
func (a *Adapter) URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
```

## Middleware

`NamedRouter` middleware is translated to the corresponding Chi middleware
context:

* `Use()` → global Chi middleware
* `Subrouter()` → Chi subrouter middleware
* route middleware → route-specific Chi middleware

## Static files

Static files are served using Chi together with `net/http`:

```go
router.Static("/static/*", "./static")
```

## Route walking

`Walk()` delegates to `chi.Walk` and prints registered routes to standard
output.

## Compatibility

NamedRouter: 0.x
Chi: v5
Go: 1.25+

This adapter is stable and follows Semantic Versioning starting with v1.0.0.

## Related projects
* [NamedRouter](https://github.com/mafalt/namedrouter)
* [Chi](https://github.com/go-chi/chi)

For the complete NamedRouter API and architecture, see the
[NamedRouter documentation](https://github.com/mafalt/namedrouter).

## License

NamedRouter Chi Adapter is released under the MIT License.
