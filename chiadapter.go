package chiadapter

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mafalt/namedrouter"
)

// chiParameterParser is a struct that implements the namedrouter.ParameterParser interface.
type chiParameterParser struct{}

// Parse parses the route pattern and returns the names of the parameters.
func (p *chiParameterParser) Parse(route namedrouter.RouteDefinition) (namedrouter.RouteParameterNames, error) {
	params := make(namedrouter.RouteParameterNames)

	for _, routePart := range strings.Split(route.Pattern, "/") {
		switch {
		case strings.HasPrefix(routePart, "{") && strings.HasSuffix(routePart, "}"):
			paramName := strings.TrimSuffix(strings.TrimPrefix(routePart, "{"), "}")
			params[paramName] = struct{}{}
		case strings.Contains(routePart, "{") || strings.Contains(routePart, "}"):
			return nil, &namedrouter.InvalidRouteParameterFormatError{RouteName: route.Name}
		}
	}

	return params, nil
}

// chiParameterApplier is a struct that implements the namedrouter.ParameterApplier interface.
type chiParameterApplier struct{}

// Apply applies the parameter values to the specified pattern.
func (a *chiParameterApplier) Apply(pattern string, params namedrouter.RouteParams) string {
	for paramName := range params {
		value := url.PathEscape(fmt.Sprint(params[paramName]))
		pattern = strings.ReplaceAll(pattern, "{"+paramName+"}", value)
	}

	return pattern
}

// chiAdapter is an implementation of the namedrouter.Adapter interface using the Chi router.
type chiAdapter struct {
	router  chi.Router
	prefix  string
	parser  namedrouter.ParameterParser
	applier namedrouter.ParameterApplier
}

// New creates and returns a new instance of chiAdapter, which implements the namedrouter.Adapter interface.
func New() namedrouter.Adapter {
	return &chiAdapter{
		router:  chi.NewRouter(),
		parser:  &chiParameterParser{},
		applier: &chiParameterApplier{},
	}
}

// Subrouter creates a new subrouter with the specified prefix and optional middlewares.
func (c *chiAdapter) Subrouter(prefix string, middlewares ...namedrouter.Middleware) namedrouter.Adapter {
	childRouter := c.router.Group(func(r chi.Router) {
		r.Use(middlewares...)
	})

	return &chiAdapter{
		router: childRouter,
		prefix: c.JoinPath(c.prefix, prefix),
	}
}

// ApplyMiddlewares applies the specified middlewares to the router.
func (c *chiAdapter) ApplyMiddlewares(pattern string, middlewares ...namedrouter.Middleware) namedrouter.Adapter {
	c.router = c.router.With(middlewares...)
	return c
}

// Register registers a new route with the specified HTTP method, pattern, name, handler, and optional middlewares.
func (c *chiAdapter) Register(method, pattern string, handler http.HandlerFunc) {
	c.router.Method(method, pattern, handler)
}

// ServeHTTP implements the http.Handler interface for the chiAdapter.
func (c *chiAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.router.ServeHTTP(w, r)
}

// Walk prints all registered routes in the router to the console.
func (c *chiAdapter) Walk() {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/")
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(c.router, walkFunc); err != nil {
		panic(err)
	}
}

// JoinPath joins the prefix and pattern paths, ensuring that the resulting path is properly formatted.
func (c *chiAdapter) JoinPath(prefix, pattern string) string {
	if prefix == "" && pattern == "" {
		return "/"
	}

	return path.Join(prefix, pattern)
}

// Use applies the specified middlewares to the router.
func (c *chiAdapter) Use(middlewares ...namedrouter.Middleware) {
	c.router.Use(middlewares...)
}

// Static serves static files from the specified root directory under the given prefix.
func (c *chiAdapter) Static(pattern, root string) {
	stripPrefix := strings.TrimSuffix(pattern, "*")
	c.router.Handle(pattern, http.StripPrefix(stripPrefix, http.FileServer(http.Dir(root))))
}

// ParameterParser returns the ParameterParser used by the chiAdapter.
func (c *chiAdapter) ParameterParser() namedrouter.ParameterParser {
	return c.parser
}

// ParameterApplier returns the ParameterApplier used by the chiAdapter.
func (c *chiAdapter) ParameterApplier() namedrouter.ParameterApplier {
	return c.applier
}

// URLParam retrieves the value of a URL parameter from the request using the Chi router.
func (c *chiAdapter) URLParam(r *http.Request, key string) string {
	if r == nil || key == "" {
		return ""
	}
	return chi.URLParam(r, key)
}
