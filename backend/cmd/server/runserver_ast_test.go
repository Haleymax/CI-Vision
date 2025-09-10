// Package server tests
// Testing framework: Go standard library "testing" (no external dependencies).
// Notes:
// - Focused on validating the structure and intent of RunServer as given in server_test.go.
// - We avoid actually starting an HTTP server to keep tests deterministic and fast.
// - Where direct runtime validation is impractical (due to log.Fatal/os.Exit and blocking server),
//   we use AST-based assertions and behavioral checks on the port conversion pattern used in the diff.
package server

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"testing"
	"unicode/utf8"
)

// getFuncDecl parses the given file and returns the *ast.FuncDecl for the named function.
func getFuncDecl(t *testing.T, filename, funcName string) *ast.FuncDecl {
	t.Helper()
	fset := token.NewFileSet()
	// Parse the file from disk relative to this package directory.
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err \!= nil {
		t.Fatalf("failed to parse %s: %v", filename, err)
	}
	for _, d := range f.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok && fn.Name \!= nil && fn.Name.Name == funcName {
			return fn
		}
	}
	return nil
}

// Test_RunServer_AST_Structure verifies that RunServer:
// - creates a Gin engine via gin.Default()
// - wires routes using routers.SetupRouter(r)
// - starts the server in a goroutine (go func(){ ... }())
// - attempts to run the engine via r.Run(...)
// - performs port conversion using string(rune(...)) as seen in the diff
func Test_RunServer_AST_Structure(t *testing.T) {
	const target = "server_test.go" // as provided in this PR context
	fd := getFuncDecl(t, target, "RunServer")
	if fd == nil {
		t.Fatalf("RunServer not found in %s", target)
	}

	var (
		foundGinDefault        bool
		foundSetupRouter       bool
		foundSetupRouterArgR   bool
		foundGoStmt            bool
		foundRun               bool
		foundStringRuneConvert bool
	)

	ast.Inspect(fd, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.GoStmt:
			foundGoStmt = true
		case *ast.CallExpr:
			switch fn := n.Fun.(type) {
			case *ast.SelectorExpr:
				// X.Sel form
				if ident, ok := fn.X.(*ast.Ident); ok {
					switch {
					case ident.Name == "gin" && fn.Sel.Name == "Default":
						foundGinDefault = true
					case ident.Name == "routers" && fn.Sel.Name == "SetupRouter":
						foundSetupRouter = true
						if len(n.Args) == 1 {
							if argIdent, ok := n.Args[0].(*ast.Ident); ok && argIdent.Name == "r" {
								foundSetupRouterArgR = true
							}
						}
					case ident.Name == "r" && fn.Sel.Name == "Run":
						foundRun = true
					}
				}
			case *ast.Ident:
				// Look for string(rune(...)) usage anywhere within RunServer
				if fn.Name == "string" && len(n.Args) == 1 {
					if inner, ok := n.Args[0].(*ast.CallExpr); ok {
						if innerIdent, ok := inner.Fun.(*ast.Ident); ok && innerIdent.Name == "rune" {
							foundStringRuneConvert = true
						}
					}
				}
			}
		}
		return true
	})

	if \!foundGinDefault {
		t.Errorf("expected gin.Default() to be called in RunServer")
	}
	if \!foundSetupRouter {
		t.Errorf("expected routers.SetupRouter(...) to be called in RunServer")
	}
	if \!foundSetupRouterArgR {
		t.Errorf("expected routers.SetupRouter to receive the gin engine variable 'r'")
	}
	if \!foundGoStmt {
		t.Errorf("expected RunServer to start the server in a goroutine (go func(){...}())")
	}
	if \!foundRun {
		t.Errorf("expected r.Run(...) to be called inside RunServer")
	}
	if \!foundStringRuneConvert {
		t.Errorf("expected use of string(rune(...)) for port conversion (as per diff); consider replacing with strconv.Itoa for correctness")
	}
}

// Test_PortConversion_UsingRune_ProducesNonNumericPort documents the behavior of using string(rune(port))
// seen in RunServer's address construction. This highlights why ports should be converted with strconv.Itoa.
func Test_PortConversion_UsingRune_ProducesNonNumericPort(t *testing.T) {
	ports := []int{0, 80, 443, 8080, 65535}
	for _, p := range ports {
		p := p
		t.Run(strconv.Itoa(p), func(t *testing.T) {
			t.Parallel()
			r := string(rune(p))
			// Expect exactly one rune (which may be multiple bytes)
			if got := utf8.RuneCountInString(r); got \!= 1 {
				t.Errorf("string(rune(%d)) produced %d runes, want 1", p, got)
			}
			// It should NOT equal the decimal string of the port (demonstrates the bug)
			if r == strconv.Itoa(p) {
				t.Errorf("string(rune(%d)) unexpectedly equals its decimal string representation", p)
			}
		})
	}
}