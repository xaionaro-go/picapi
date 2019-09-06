package config

import "time"

// Config is used to configure `picapi`, see fields below.
type Config struct {
	// ListenAddress defines the address to be used to listen for incoming HTTP requests
	//
	// Default value: ":8486"
	ListenAddress string `envconfig:"LISTEN_ADDRESS" default:":8486"`

	// LoggingLevel defines logging level of the application.
	//
	// Allowed values: "fatal", "debug"
	//
	// Default value: "fatal"
	LoggingLevel string `envconfig:"LOGGING_LEVEL" default:"fatal"`

	// CacheDuration defines how long a response is considered to be actual.
	//
	// A response couldn't be returned from the cache if it was generated
	// more than CacheDuration ago.
	CacheDuration time.Duration `envconfig:"CACHE_DURATION" default:"1h"`

	// CacheMaxEntries defines how many entries could be placed into the cache.
	//
	// If you increase this value then the service may consume more memory.
	//
	// Total memory consumption by the cache is limited by CacheMaxEntries * CacheMaxEntrySize.
	CacheMaxEntries uint64 `envconfig:"CACHE_MAX_ENTRIES" default:"65536"`

	// CacheMaxEntrySize defines how large entries could be placed into the cache.
	//
	// If you increase this value then the service may consume more memory.
	//
	// Total memory consumption by the cache is limited by CacheMaxEntries * CacheMaxEntrySize.
	CacheMaxEntrySize uint64 `envconfig:"CACHE_MAX_ENTRY_SIZE" default:"1048576"`
}
