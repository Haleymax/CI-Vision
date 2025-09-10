// Note: Tests use Go's standard library "testing" package.

package main

import (
	"sync/atomic"
	"testing"
)

// Test that main calls runServer exactly once (happy path).
func TestMainInvokesRunServerOnce(t *testing.T) {
	var called int32
	prev := runServer
	t.Cleanup(func() { runServer = prev })

	runServer = func() { atomic.AddInt32(&called, 1) }

	main()

	if got := atomic.LoadInt32(&called); got \!= 1 {
		t.Fatalf("expected runServer to be called once; got %d", got)
	}
}

// Test that any panic from runServer propagates up through main (failure condition).
func TestMainPropagatesPanicFromRunServer(t *testing.T) {
	prev := runServer
	t.Cleanup(func() { runServer = prev })

	runServer = func() { panic("boom") }

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from runServer to propagate; got none")
		}
	}()

	main()
}

// Test that multiple invocations of main() call the current runServer reference each time.
func TestMainCanBeCalledMultipleTimes(t *testing.T) {
	var called int32
	prev := runServer
	t.Cleanup(func() { runServer = prev })

	runServer = func() { atomic.AddInt32(&called, 1) }

	main()
	main()

	if got := atomic.LoadInt32(&called); got \!= 2 {
		t.Fatalf("expected runServer to be called twice; got %d", got)
	}
}