package httpserver

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/xaionaro-go/picapi/imageprocessor"
)

const (
	maxWidth  = 5000
	maxHeight = 5000
)

var (
	httpClient = &http.Client{}
)

// getParamsResize parses and validates GET-arguments for method `GET /resize`
func (srv *HTTPServer) getParamsResize(ctx *fasthttp.RequestCtx) (width, height uint, sourceURL *url.URL, err error) {
	// Parsing "width"
	{
		widthString := string(ctx.QueryArgs().Peek(`width`))
		if widthString == `` {
			err = &ErrBadRequest{`"width" is cannot be empty`}
			return
		}

		var width64 uint64
		width64, err = strconv.ParseUint(widthString, 10, 64)
		if err != nil {
			err = &ErrBadRequest{`"width" is invalid: ` + err.Error()}
			return
		}

		if width64 < 1 {
			err = &ErrBadRequest{`"width" should be positive`}
			return
		}

		if width64 > maxWidth {
			err = &ErrBadRequest{`"width" is too high`}
			return
		}

		width = uint(width64)
	}

	// Parsing "height"
	{
		heightString := string(ctx.QueryArgs().Peek(`height`))
		if heightString == `` {
			err = &ErrBadRequest{`"height" is cannot be empty`}
			return
		}

		var height64 uint64
		height64, err = strconv.ParseUint(heightString, 10, 64)
		if err != nil {
			err = &ErrBadRequest{`"height" is invalid: ` + err.Error()}
			return
		}

		if height64 < 1 {
			err = &ErrBadRequest{`"height" should be positive`}
			return
		}

		if height64 > maxHeight {
			err = &ErrBadRequest{`"height" is too high`}
			return
		}

		height = uint(height64)
	}

	// Parsing "url"
	{
		sourceURLString := string(ctx.QueryArgs().Peek(`url`))
		if sourceURLString == `` {
			err = &ErrBadRequest{`"url" is cannot be empty`}
			return
		}

		sourceURL, err = url.Parse(sourceURLString)
		if err != nil {
			err = &ErrBadRequest{`"url" is invalid: ` + err.Error()}
			return
		}

		if sourceURLString != `test_picture` && len(sourceURL.Scheme) == 0 {
			err = &ErrBadRequest{`"url" is invalid: the scheme cannot be empty`}
			return
		}
	}

	return
}

// handleResize is the handler for `GET /resize`
func (srv *HTTPServer) handleResize(ctx *fasthttp.RequestCtx) {
	// Caution:
	// * This method could be used as a DDoS-attack multiplier. So it's required to hide it from the public access.
	// * There's no stoppers if the server is already overloaded (or there're obviously too many requests).
	// * Here's pretty stupid HTTP-getter, so there could be a lot of hanged connections.
	// * We don't check a resolution of the incoming image

	// Initializing

	startTS := time.Now()
	defer func() {
		ctx.Response.Header.Set(`X-Profiling-Value-Total`, time.Now().Sub(startTS).String())
	}()

	width, height, sourceURL, err := srv.getParamsResize(ctx)
	if err != nil {
		srv.handleError(ctx, err)
		return
	}

	// Download the image

	// TODO: consider if we should cache downloaded images

	downloadStartTS := time.Now()
	sourceBody, err := srv.download(ctx.Done(), `GET`, sourceURL)
	if err != nil {
		srv.handleError(ctx, err)
		return
	}
	defer sourceBody.Close()
	ctx.Response.Header.Set(`X-Profiling-Value-Download`, time.Now().Sub(downloadStartTS).String())

	// Process the image

	sourceReader := sourceBody

	imageFormat, err := srv.ImageProcessor.Resize(
		sourceReader,
		ctx.Response.BodyWriter(),
		width,
		height,
	)
	switch imageFormat {
	case imageprocessor.ImageFormatJPEG:
		ctx.Response.Header.Set(`Content-Type`, `image/jpeg`)
	default:
		err = &ErrBadRequest{fmt.Sprintf(`unexpected image format "%v" (expected: "JPEG")`, imageFormat)}
	}
	if err != nil {
		srv.handleError(ctx, err)
		return
	}

	return
}
