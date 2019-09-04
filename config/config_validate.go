package config

import (
	"fmt"
	"strings"

	"github.com/xaionaro-go/errors"
)

// Validate checks if the configuration is correct and returns an error if is not
func (cfg *Config) Validate() (err error) {
	defer func() { err = errors.Wrap(err) }()

	validLoggingLevels := []string{`fatal`, `debug`}
	isValidLoggingLevel := false
	for _, validLogginLevel := range validLoggingLevels {
		if cfg.LoggingLevel == validLogginLevel {
			isValidLoggingLevel = true
			break
		}
	}

	if !isValidLoggingLevel {
		return fmt.Errorf(`invalid logging level: %v (valid values: %v)`,
			cfg.LoggingLevel,
			strings.Join(validLoggingLevels, `,`),
		)
	}

	return nil
}
