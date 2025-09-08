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
