package chiadapter_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mafalt/namedrouter"
	chiadapter "github.com/mafalt/namedrouter-chi"
)

func TestChiParameterParser_Parse(t *testing.T) {
	parser := chiadapter.New().ParameterParser()

	params, err := parser.Parse(namedrouter.RouteDefinition{
		Name:    "tenants.users.show",
		Pattern: "/tenants/{tenantId}/users/{userId}",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !params.Contains("tenantId") {
		t.Fatal("Expected tenantId parameter to be parsed")
	}
	if !params.Contains("userId") {
		t.Fatal("Expected userId parameter to be parsed")
	}
}

func TestChiParameterParser_Parse_InvalidFormat(t *testing.T) {
	parser := chiadapter.New().ParameterParser()

	_, err := parser.Parse(namedrouter.RouteDefinition{
		Name:    "users.show",
		Pattern: "/users/{id",
	})
	if err == nil {
		t.Fatal("Expected error for invalid parameter format")
	}

	if _, ok := err.(*namedrouter.InvalidRouteParameterFormatError); !ok {
		t.Fatalf("Expected InvalidRouteParameterFormatError, got %T", err)
	}
}

func TestChiParameterApplier_Apply_EncodesValues(t *testing.T) {
	applier := chiadapter.New().ParameterApplier()

	url := applier.Apply("/users/{name}", namedrouter.RouteParams{"name": "john doe"})
	if url != "/users/john%20doe" {
		t.Fatalf("Expected URL to be encoded, got %s", url)
	}
}

func TestChiAdapter_JoinPath(t *testing.T) {
	adapter := chiadapter.New()

	if got := adapter.JoinPath("", ""); got != "/" {
		t.Fatalf("Expected / for empty path, got %s", got)
	}
	if got := adapter.JoinPath("/api", "/v1"); got != "/api/v1" {
		t.Fatalf("Expected /api/v1, got %s", got)
	}
}

func TestChiAdapter_URLParam(t *testing.T) {
	adapter := chiadapter.New()

	t.Run("nil request returns empty string", func(t *testing.T) {
		if got := adapter.URLParam(nil, "id"); got != "" {
			t.Fatalf("Expected empty string from nil request, got %q", got)
		}
	})

	t.Run("empty key returns empty string", func(t *testing.T) {
		if got := adapter.URLParam(httptest.NewRequest(http.MethodGet, "/users/42", nil), ""); got != "" {
			t.Fatalf("Expected empty string for empty key, got %q", got)
		}
	})

	t.Run("returns matched route parameter", func(t *testing.T) {
		router := chi.NewRouter()
		router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			if got := adapter.URLParam(r, "id"); got != "42" {
				t.Fatalf("Expected route param id to be 42, got %q", got)
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("missing key returns empty string", func(t *testing.T) {
		router := chi.NewRouter()
		router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			if got := adapter.URLParam(r, "missing"); got != "" {
				t.Fatalf("Expected missing key to return empty string, got %q", got)
			}
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", w.Code)
		}
	})
}

func TestChiAdapter_ApplyMiddlewares(t *testing.T) {
	adapter := chiadapter.New()
	middlewareCalled := false

	adapterWithMiddleware := adapter.ApplyMiddlewares("/health", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalled = true
			next.ServeHTTP(w, r)
		})
	})
	adapterWithMiddleware.Register(http.MethodGet, "/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	adapter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	if !middlewareCalled {
		t.Fatal("Expected middleware to be called")
	}
}

func TestChiAdapter_Static_UsesProvidedPattern(t *testing.T) {
	adapter := chiadapter.New()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "hello.txt")

	if err := os.WriteFile(filePath, []byte("hello"), 0o600); err != nil {
		t.Fatalf("Failed to create static test file: %v", err)
	}

	adapter.Static("/assets/*", tempDir)

	req := httptest.NewRequest(http.MethodGet, "/assets/hello.txt", nil)
	w := httptest.NewRecorder()
	adapter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "hello" {
		t.Fatalf("Expected body hello, got %s", w.Body.String())
	}
}
