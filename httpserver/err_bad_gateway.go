package httpserver

// ErrBadGateway is an error which will be returned to a browser with as
// a "Bad Gateway" error
type ErrBadGateway struct {
	Message string
}

// Error just implements interface `error`.
func (err *ErrBadGateway) Error() string {
	return `bad gateway: ` + err.Message
}
