package config

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
}
