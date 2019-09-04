package httpserver

import (
	"time"

	"github.com/buaazp/fasthttprouter"

	"github.com/xaionaro-go/picapi/httpserver/middlewares"
)

func (srv *HTTPServer) newRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.GET(`/resize`, middlewares.Cache(time.Hour, srv.handleResize))

	return router
}
