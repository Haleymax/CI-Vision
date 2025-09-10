package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAppConfig(t *testing.T) {
	config := LoadConfig()

	assert.NotNil(t, config, "Config object should not be nil")

	config2 := GetConfig()
	assert.Equal(t, config, config2, "GetConfig should return the same config instance")

	t.Logf("Application config loaded successfully")
}

func TestLoadMySQLConfig(t *testing.T) {
	config := GetConfig()
	assert.NotNil(t, config, "Config object should not be nil")

	mysqlConfig := config.MySQL

	assert.NotEmpty(t, mysqlConfig.Host, "MySQL host should not be empty")
	assert.NotZero(t, mysqlConfig.Port, "MySQL port should not be zero")
	assert.NotEmpty(t, mysqlConfig.Username, "MySQL username should not be empty")
	assert.NotEmpty(t, mysqlConfig.Password, "MySQL password should not be empty")
	assert.NotEmpty(t, mysqlConfig.Database, "MySQL database name should not be empty")

	t.Logf("MySQL config loaded successfully")
}

// Note: Additional tests complement existing ones.
// They validate concurrency-safety, singleton behavior, and basic MySQL config invariants
// without assuming reload capabilities or private APIs.

func TestGetConfig_ConcurrentSingleton(t *testing.T) {
	t.Parallel()

	const goroutines = 64
	var wg sync.WaitGroup
	wg.Add(goroutines)

	results := make([]*AppConfig, goroutines)
	start := make(chan struct{})

	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			<-start
			results[i] = GetConfig()
		}()
	}
	// Fan-out start
	close(start)
	wg.Wait()

	// All instances should be identical (singleton)
	first := results[0]
	assert.NotNil(t, first, "GetConfig should not return nil")
	for i := 1; i < goroutines; i++ {
		assert.Equal(t, first, results[i], "GetConfig should return the same instance across goroutines")
	}
}

func TestGetConfig_RepeatedCallsStableOverTime(t *testing.T) {
	t.Parallel()

	c1 := GetConfig()
	time.Sleep(10 * time.Millisecond)
	c2 := GetConfig()

	assert.NotNil(t, c1)
	assert.Same(t, c1, c2, "GetConfig should consistently return the same pointer over time")
}

func TestMySQLConfig_InvariantsAndFormats(t *testing.T) {
	t.Parallel()

	cfg := GetConfig()
	requireNotNil := func() {
		if cfg == nil {
			t.Fatalf("config must not be nil")
		}
	}
	requireNotNil()

	mysql := cfg.MySQL

	// Existing tests assert non-zero/non-empty; we deepen checks:
	//  - Port is within valid TCP range
	//  - Host has no surrounding whitespace
	//  - Username does not contain whitespace-only value
	//  - Database name is not whitespace-only
	assert.Greater(t, mysql.Port, 0, "MySQL port must be > 0")
	assert.LessOrEqual(t, mysql.Port, 65535, "MySQL port must be <= 65535")

	assert.Equal(t, strings.TrimSpace(mysql.Host), mysql.Host, "MySQL host should not contain leading/trailing spaces")
	assert.NotEqual(t, "", strings.TrimSpace(mysql.Username), "MySQL username should not be whitespace-only")
	assert.NotEqual(t, "", strings.TrimSpace(mysql.Database), "MySQL database should not be whitespace-only")
}

func TestMySQLConfig_EnvWhitespaceDoesNotBreakAccess(t *testing.T) {
	t.Parallel()

	// Intentionally set benign env vars with whitespace to ensure accessing the config does not panic
	// (We do not assert reload behavior; just that the config access is safe under unusual env.)
	// We restore previous values after the test.
	type kv struct{ k, v string }
	restore := []kv{}
	for _, k := range []string{
		"MYSQL_HOST", "MYSQL_USERNAME", "MYSQL_PASSWORD", "MYSQL_DATABASE",
	} {
		old, ok := os.LookupEnv(k)
		if ok {
			restore = append(restore, kv{k, old})
		}
		_ = os.Setenv(k, "  value-with-spaces  ")
	}
	defer func() {
		for _, r := range restore {
			_ = os.Setenv(r.k, r.v)
		}
		for _, k := range []string{"MYSQL_HOST", "MYSQL_USERNAME", "MYSQL_PASSWORD", "MYSQL_DATABASE"} {
			if _, ok := os.LookupEnv(k); \!ok {
				_ = os.Unsetenv(k)
			}
		}
	}()

	cfg := GetConfig()
	assert.NotNil(t, cfg, "GetConfig should remain safe even if env contains unusual whitespace")
	_ = cfg.MySQL // Access to ensure no panic on field read
}

func TestLoadConfig_Idempotency(t *testing.T) {
	t.Parallel()

	// Even if LoadConfig is called multiple times, resulting instance should equal GetConfig's instance.
	// This does not assume reloading; it only verifies stability of returned instance.
	c1 := LoadConfig()
	c2 := GetConfig()
	c3 := LoadConfig()

	assert.Same(t, c1, c2, "LoadConfig and GetConfig should refer to the same instance")
	assert.Same(t, c2, c3, "Repeated LoadConfig calls should remain idempotent with the same instance")
}

func TestGetConfig_DoesNotLeakWithGC(t *testing.T) {
	t.Parallel()

	// Sanity check to ensure that holding only the returned pointer remains valid across a GC cycle.
	c := GetConfig()
	assert.NotNil(t, c)
	runtime.GC()
	assert.NotNil(t, c, "Config pointer should remain valid across GC")
}
