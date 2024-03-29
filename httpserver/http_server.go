package httpserver

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/xaionaro-go/errors"
	"github.com/xaionaro-go/picapi/httpserver/middlewares"
	"github.com/xaionaro-go/picapi/imageprocessor"
)

const (
	version = `0.1`
)

var (
	// ErrAlreadyStarted is returned when method `Start` is called on an already started `HTTPServer`.
	ErrAlreadyStarted = errors.New(`already started`)

	// ErrNotStarted is returned when method `Stop` is called on a non-started-yet `HTTPServer`.
	ErrNotStarted = errors.New(`not started`)

	// ErrImageProcessorIsRequireed is returned when NewHTTPServer is called with nil `proc`.
	ErrImageProcessorIsRequireed = errors.New(`image processor is required`)
)

// HTTPServer is the implemention of an HTTP server which implements `picapi`'s methods
type HTTPServer struct {
	startStopLocker sync.Mutex

	context       context.Context
	stopFunc      context.CancelFunc
	stopWaitGroup sync.WaitGroup

	pendingError error

	httpListener net.Listener
	httpBackend  *fasthttp.Server
	httpRouter   *fasthttprouter.Router

	accessLogger Printfer

	// ImageProcessor is the implementation of image manipulation tool
	ImageProcessor imageprocessor.ImageProcessor
}

// Printfer is a simple logger interface for `NewHTTPServer`
type Printfer interface {
	// Printf is a pretty-standart `Printf` (`fmt` is the format and `args` are the arguments for the format).
	//
	// For example you can use `log.New` as `Printfer`.
	Printf(fmt string, args ...interface{})
}

// NewHTTPServer returns a new instance of `HTTPServer`.
//
// `HTTPServer` will not be started automatically. It's required to call method `Start`
// to start it.
//
// `accessLogger` is used to log requests (similar to an "access.log").
//
// `handlerLogger` is used to log errors and warnings by handlers (similar to "error.log")
func NewHTTPServer(
	proc imageprocessor.ImageProcessor,
	accessLogger Printfer,
	handlerLogger Printfer,
	cacheDuration time.Duration,
	cacheMaxEntries uint64,
	cacheMaxEntrySize uint64,
) (srv *HTTPServer, err error) {
	defer func() { err = errors.Wrap(err) }()

	if proc == nil {
		return nil, ErrImageProcessorIsRequireed
	}

	srv = &HTTPServer{
		ImageProcessor: proc,
	}

	srv.httpBackend = &fasthttp.Server{
		Name:   `picapi ` + version,
		Logger: handlerLogger, // TODO: reconsider logging system to correctly comply to logging levels
	}
	srv.httpRouter = srv.newRouter(
		cacheDuration,
		cacheMaxEntries,
		cacheMaxEntrySize,
	)

	handler := middlewares.RecoverPanic(srv.httpRouter.Handler)
	if accessLogger != nil {
		srv.accessLogger = accessLogger
		handler = middlewares.AccessLogger(accessLogger, handler)
	}
	srv.httpBackend.Handler = handler
	return
}

// Start creates a listener and starts the process of serving incoming HTTP requests.
func (srv *HTTPServer) Start(
	ctx context.Context,
	listenAddress string,
) (err error) {
	defer func() { err = errors.Wrap(err) }()

	srv.startStopLocker.Lock()
	defer srv.startStopLocker.Unlock()

	if srv.context != nil {
		return ErrAlreadyStarted
	}

	srv.httpListener, err = net.Listen(`tcp`, listenAddress)
	if err != nil {
		return
	}

	if srv.accessLogger != nil {
		srv.accessLogger.Printf("started to listen at %v\n", listenAddress)
	}

	srv.context, srv.stopFunc = context.WithCancel(ctx)
	srv.stopWaitGroup.Add(2)
	go func() {
		defer srv.stopWaitGroup.Done()

		srv.pendingError = errors.Wrap(
			srv.httpBackend.Serve(srv.httpListener),
		)
	}()
	go func() {
		defer srv.stopWaitGroup.Done()

		select {
		case <-srv.Done():
		}
		if srv.accessLogger != nil {
			srv.accessLogger.Printf("stopping\n")
		}
		srv.httpBackend.Shutdown()
		if srv.accessLogger != nil {
			srv.accessLogger.Printf("stopped\n")
		}
		srv.httpListener.Close()
		srv.context = nil
	}()

	runtime.Gosched() // to allow `Serve` to initialize the server
	return
}

// Stop is anti-Start: it stops the process of serving incoming HTTP requsts and closes
// the listener.
func (srv *HTTPServer) Stop() (err error) {
	defer func() { err = errors.Wrap(err) }()

	srv.startStopLocker.Lock()
	defer srv.startStopLocker.Unlock()

	if srv.context == nil {
		return ErrNotStarted
	}

	srv.stopFunc()
	srv.Wait()
	return nil
}

// Done returns the `Done()` of the context of the HTTPServer instance.
func (srv *HTTPServer) Done() <-chan struct{} {
	return srv.context.Done()
}

// Wait waits for HTTP server to stop
func (srv *HTTPServer) Wait() error {
	srv.stopWaitGroup.Wait()
	return srv.pendingError
}

// handleError is a helper to handle error cases from real handlers (like `handleResize`).
func (srv *HTTPServer) handleError(ctx *fasthttp.RequestCtx, err error) {
	var status int

	ctx.Response.ResetBody()

	switch err := err.(type) {
	case *ErrBadRequest:
		status = http.StatusBadRequest
	case *ErrBadGateway:
		status = http.StatusBadGateway
	case *ErrForbidden:
		status = http.StatusForbidden
	default:
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Logger().Printf(`[handleError] unknown error: %v`, err)
		return
	}

	ctx.Response.SetStatusCode(status)
	ctx.Response.BodyWriter().Write([]byte(err.Error()))
	return
}

// download is a helper to download files from a remote parties
//
// Special value of sourceURL is `test_picture`: it will return a static 10x10 picture (with noise on it) from memory.
func (srv *HTTPServer) download(cancelChan <-chan struct{}, method string, sourceURL *url.URL) (io.ReadCloser, error) {
	if method == `GET` && sourceURL.String() == `test_picture` {
		return testPictureBody(), nil
	}

	sourceRequest, err := http.NewRequest(method, sourceURL.String(), nil)
	if err != nil {
		return nil, &ErrBadRequest{`"url" is invalid: unable to prepare a request: ` + err.Error()}
	}
	sourceRequest.Cancel = cancelChan

	sourceResponse, err := httpClient.Do(sourceRequest)
	if err != nil {
		return nil, &ErrBadGateway{err.Error()}
	}
	return sourceResponse.Body, nil
}
