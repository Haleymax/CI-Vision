package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRouter_RegistersHelloRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	SetupRouter(engine)

	// Attempt to hit the expected hello endpoint under /api.
	// If the exact path differs, adjust to match the groups.HelloRouters definition.
	req := httptest.NewRequest(http.MethodGet, "/api/hello", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code \!= http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d, body: %s", w.Code, w.Body.String())
	}
	// Basic body assertion; adapt if the handler returns JSON.
	if len(w.Body.Bytes()) == 0 {
		t.Fatalf("expected non-empty response body")
	}
}

func TestSetupRouter_UnknownRouteUnderAPI_Returns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	SetupRouter(engine)

	req := httptest.NewRequest(http.MethodGet, "/api/unknown-route", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if w.Code \!= http.StatusNotFound && w.Code \!= http.StatusMethodNotAllowed {
		t.Fatalf("expected 404 Not Found or 405 Method Not Allowed, got %d", w.Code)
	}
}

func TestSetupRouter_PanicsOnNilEngine(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when passing nil engine to SetupRouter")
		}
	}()
	SetupRouter(nil)
}

func TestSetupRouter_GroupPrefixDoesNotMatchRoot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	SetupRouter(engine)

	// Root path should not accidentally match API routes.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// Typically 404 for root unless other routes are registered elsewhere.
	if w.Code == http.StatusOK {
		t.Fatalf("expected non-200 at '/', got 200; API routes should be under /api/*")
	}
}