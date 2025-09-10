// Tests use Go standard library "testing" with github.com/stretchr/testify/require for assertions.
// Focus: Validate NewControllers behavior added in setup/controllers.go, especially construction of HelloController value.

package setup

import (
	"reflect"
	"sync"
	"testing"

	"civ/internal/controller/hello"

	"github.com/stretchr/testify/require"
)

// TestNewControllers_ReturnsNonNil ensures constructor returns a non-nil pointer.
func TestNewControllers_ReturnsNonNil(t *testing.T) {
	t.Parallel()
	got := NewControllers()
	require.NotNil(t, got, "NewControllers should not return nil")
}

// TestNewControllers_HelloControllerTypeAndValue verifies HelloController has expected type
// and equals a freshly constructed hello.HelloController value (from hello.NewHelloController()).
func TestNewControllers_HelloControllerTypeAndValue(t *testing.T) {
	t.Parallel()
	got := NewControllers()
	require.NotNil(t, got)

	ref := hello.NewHelloController()
	require.NotNil(t, ref, "hello.NewHelloController should not return nil")

	// Type check
	gotType := reflect.TypeOf(got.HelloController)
	wantType := reflect.TypeOf(*ref)
	require.Equal(t, wantType, gotType, "HelloController concrete type mismatch")

	// Value equality (constructor currently returns zero-value struct)
	require.Equal(t, *ref, got.HelloController, "HelloController value should match a fresh constructor value")
}

// TestNewControllers_ReturnsFreshInstances ensures each call returns a distinct *Controllers pointer
// and that HelloController fields live at distinct addresses (no aliasing across instances).
func TestNewControllers_ReturnsFreshInstances(t *testing.T) {
	t.Parallel()
	a := NewControllers()
	b := NewControllers()

	require.NotNil(t, a)
	require.NotNil(t, b)

	// Distinct *Controllers instances
	require.NotEqual(t, a, b, "Expected distinct *Controllers instances from separate calls")

	// Distinct HelloController field addresses
	require.NotEqual(t, &a.HelloController, &b.HelloController, "HelloController fields should not alias across instances")
}

// TestNewControllers_ConcurrentCalls does a light concurrency sanity check.
func TestNewControllers_ConcurrentCalls(t *testing.T) {
	t.Parallel()
	const N = 16

	var wg sync.WaitGroup
	wg.Add(N)

	results := make(chan *Controllers, N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			results <- NewControllers()
		}()
	}

	wg.Wait()
	close(results)

	seen := map[*Controllers]struct{}{}
	seenHC := map[*hello.HelloController]struct{}{}

	for c := range results {
		require.NotNil(t, c, "NewControllers returned nil in concurrent call")

		// Ensure unique *Controllers pointers
		if _, dup := seen[c]; dup {
			t.Fatalf("duplicate *Controllers pointer returned across concurrent calls: %p", c)
		}
		seen[c] = struct{}{}

		// Ensure HelloController field addresses are unique across instances
		hcPtr := &c.HelloController
		if _, dup := seenHC[hcPtr]; dup {
			t.Fatalf("duplicate HelloController field address across instances: %p", hcPtr)
		}
		seenHC[hcPtr] = struct{}{}
	}
}