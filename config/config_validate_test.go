package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidateLoggingLevel(t *testing.T) {
	cfg := &Config{}

	cfg.LoggingLevel = `debug`
	assert.NoError(t, cfg.Validate())

	cfg.LoggingLevel = `invalid`
	assert.Error(t, cfg.Validate())
}
