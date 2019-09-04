package httpserver

// ErrForbidden is an error which will be returned to a browser with as
// a "Forbidden" error
type ErrForbidden struct {
	Message string
}

// Error just implements interface `error`.
func (err *ErrForbidden) Error() string {
	return `forbidden: ` + err.Message
}
