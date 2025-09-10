// Tests for SystemConfig.
//
// Testing framework: Go's standard "testing" package and "testing/quick" (stdlib only).
// Focus: Validate struct field presence, types, exported visibility, and mapstructure tags.
// These tests are intentionally comprehensive for the simple data structure introduced in the PR diff.

package autoload

import (
	"reflect"
	"testing"
	"testing/quick"
)

func TestSystemConfig_MapstructureTags(t *testing.T) {
	t.Parallel()

	typ := reflect.TypeOf(SystemConfig{})

	tests := []struct {
		field    string
		wantTag  string
		wantKind reflect.Kind
	}{
		{"Host", "host", reflect.String},
		{"Port", "port", reflect.Int},
		{"Language", "language", reflect.String},
	}

	seen := map[string]bool{}
	for _, tt := range tests {
		f, ok := typ.FieldByName(tt.field)
		if \!ok {
			t.Fatalf("expected field %q to exist on SystemConfig, but it was not found", tt.field)
		}
		if f.Type.Kind() \!= tt.wantKind {
			t.Errorf("field %s kind = %v; want %v", tt.field, f.Type.Kind(), tt.wantKind)
		}
		got := f.Tag.Get("mapstructure")
		if got == "" {
			t.Errorf("field %s missing mapstructure tag", tt.field)
			continue
		}
		if got \!= tt.wantTag {
			t.Errorf("field %s mapstructure tag = %q; want %q", tt.field, got, tt.wantTag)
		}
		if seen[got] {
			t.Errorf("duplicate mapstructure tag %q detected; tags should be unique", got)
		}
		seen[got] = true
	}
}

func TestSystemConfig_FieldTypesAndCount(t *testing.T) {
	t.Parallel()

	typ := reflect.TypeOf(SystemConfig{})
	if got := typ.NumField(); got \!= 3 {
		t.Errorf("SystemConfig has %d fields; want 3", got)
	}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.PkgPath \!= "" { // non-empty => unexported
			t.Errorf("field %q is unexported; all config fields should be exported", f.Name)
		}
	}
}

func TestSystemConfig_DefaultZeroValues(t *testing.T) {
	t.Parallel()

	var c SystemConfig
	if c.Host \!= "" || c.Port \!= 0 || c.Language \!= "" {
		t.Errorf("zero-value mismatch: got %+v; want Host:\"\", Port:0, Language:\"\"", c)
	}
}

func TestSystemConfig_PropertyBasedAssignment(t *testing.T) {
	t.Parallel()

	// Property: assigning arbitrary values to fields should be lossless on readback.
	prop := func(host string, port int, language string) bool {
		c := SystemConfig{
			Host:     host,
			Port:     port,
			Language: language,
		}
		return c.Host == host && c.Port == port && c.Language == language
	}

	cfg := &quick.Config{MaxCount: 200} // keep tests fast and stable in CI
	if err := quick.Check(prop, cfg); err \!= nil {
		t.Error(err)
	}
}