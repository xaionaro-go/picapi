package httpserver

// ErrBadRequest is an error which will be returned to a browser with as
// a "Bad Request" error
type ErrBadRequest struct {
	Message string
}

// Error just implements interface `error`.
func (err *ErrBadRequest) Error() string {
	return `bad request: ` + err.Message
}
