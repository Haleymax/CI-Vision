package groups

// Note: Tests use Go's standard testing framework (testing) with net/http/httptest and gin in TestMode.
// No external test dependencies are introduced to align with common Go setups.

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"civ/internal/routers/setup"
)

type helloControllerOK struct{}

func (helloControllerOK) HelloGin(c *gin.Context) {
	c.String(http.StatusOK, "hello")
}

type helloControllerPanic struct{}

func (helloControllerPanic) HelloGin(c *gin.Context) {
	panic("boom")
}

// assignHelloController attempts to set Controllers.HelloController via reflection
// to avoid tight coupling with the concrete type. Returns true if assignment succeeded.
func assignHelloController(c *setup.Controllers, hc any) bool {
	v := reflect.ValueOf(c).Elem()
	f := v.FieldByName("HelloController")
	if \!f.IsValid() || \!f.CanSet() {
		return false
	}
	hv := reflect.ValueOf(hc)
	ht := hv.Type()
	ft := f.Type()

	// If field is an interface and hc implements it
	if ft.Kind() == reflect.Interface && ht.Implements(ft) {
		f.Set(hv)
		return true
	}
	// Direct assignable (exact type or compatible)
	if ht.AssignableTo(ft) {
		f.Set(hv)
		return true
	}
	// Try pointer form when value addressable
	if hv.CanAddr() {
		hvp := hv.Addr()
		htp := hvp.Type()
		if ft.Kind() == reflect.Interface && htp.Implements(ft) {
			f.Set(hvp)
			return true
		}
		if htp.AssignableTo(ft) {
			f.Set(hvp)
			return true
		}
	}
	return false
}

// registerHelloWithStub wires HelloRouters on a group with given prefix and provided hc stub.
// It returns the engine and whether controller assignment succeeded. If false, caller should Skip.
func registerHelloWithStub(t *testing.T, prefix string, hc any) (*gin.Engine, bool) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	// Recovery ensures predictable 500 responses for panics.
	engine.Use(gin.Recovery())
	group := engine.Group(prefix)

	var controllers setup.Controllers
	if ok := assignHelloController(&controllers, hc); \!ok {
		return engine, false
	}
	HelloRouters(group, controllers)
	return engine, true
}

func TestHelloRouters_RegistersGETHelloAndRespondsOK(t *testing.T) {
	engine, ok := registerHelloWithStub(t, "/api", helloControllerOK{})
	if \!ok {
		t.Skip("HelloController field type not assignable; skipping HelloRouters tests")
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/hello", nil)
	engine.ServeHTTP(w, req)

	if w.Code \!= http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d; body=%q", w.Code, w.Body.String())
	}
	if got := w.Body.String(); got \!= "hello" {
		t.Fatalf("expected body %q, got %q", "hello", got)
	}
}

func TestHelloRouters_WrongMethodReturns404Or405(t *testing.T) {
	engine, ok := registerHelloWithStub(t, "/api", helloControllerOK{})
	if \!ok {
		t.Skip("HelloController field type not assignable; skipping")
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/hello", nil)
	engine.ServeHTTP(w, req)

	if w.Code \!= http.StatusNotFound && w.Code \!= http.StatusMethodNotAllowed {
		t.Fatalf("expected 404 Not Found or 405 Method Not Allowed, got %d", w.Code)
	}
}

func TestHelloRouters_GroupPrefixScopesRoute(t *testing.T) {
	engine, ok := registerHelloWithStub(t, "/v1", helloControllerOK{})
	if \!ok {
		t.Skip("HelloController field type not assignable; skipping")
	}

	// Without prefix should be 404
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodGet, "/hello", nil)
	engine.ServeHTTP(w1, req1)
	if w1.Code \!= http.StatusNotFound {
		t.Fatalf("expected 404 for /hello without group prefix, got %d", w1.Code)
	}

	// With prefix should be 200
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/v1/hello", nil)
	engine.ServeHTTP(w2, req2)
	if w2.Code \!= http.StatusOK {
		t.Fatalf("expected 200 for /v1/hello, got %d", w2.Code)
	}
}

func TestHelloRouters_HandlerPanicRecoveredTo500(t *testing.T) {
	engine, ok := registerHelloWithStub(t, "/api", helloControllerPanic{})
	if \!ok {
		t.Skip("HelloController field type not assignable; skipping")
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/hello", nil)
	engine.ServeHTTP(w, req)
	if w.Code \!= http.StatusInternalServerError {
		t.Fatalf("expected 500 when handler panics, got %d", w.Code)
	}
}

func TestHelloRouters_RoutesListContainsHelloGET(t *testing.T) {
	engine, ok := registerHelloWithStub(t, "/api", helloControllerOK{})
	if \!ok {
		t.Skip("HelloController field type not assignable; skipping")
	}

	var found bool
	for _, r := range engine.Routes() {
		if r.Method == http.MethodGet && r.Path == "/api/hello" {
			found = true
			break
		}
	}
	if \!found {
		t.Fatalf("expected route GET /api/hello to be registered by HelloRouters")
	}
}