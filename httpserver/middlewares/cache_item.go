package middlewares

import (
	"time"

	"github.com/valyala/fasthttp"

	"github.com/xaionaro-go/errors"
)

var (
	// ErrCacheItemExpired is returned when it was called an `Apply()` on an expired cache item.
	ErrCacheItemExpired = errors.New(`cache item has expired`)
)

type cacheItem struct {
	// ExpireTS defines the timestamp after which the cache item is considered expired
	ExpireTS time.Time

	// Body is a copy of the HTTP response body
	Body []byte

	// ContentType is the value of HTTP header `Content-Type`
	ContentType string

	// StatusCode is the HTTP status code
	StatusCode int
}

// newCacheItem returns a new cache item with a copy of the response "resp".
//
// The item will not be valid after the moment `expireTS`.
func newCacheItem(resp *fasthttp.Response, expireTS time.Time) *cacheItem {
	c := &cacheItem{
		ExpireTS:   expireTS,
		StatusCode: resp.StatusCode(),
	}

	body := resp.Body()
	c.Body = make([]byte, len(body))
	copy(c.Body, body)

	c.ContentType = string(resp.Header.Peek(`Content-Type`))

	return c
}

// Apply sends the response from the cache.
//
// This method returns error ErrCacheItemExpired if used on an expired item.
// Also there could be other errors (for example caused by `resp.BodyWriter().Write()).
func (c *cacheItem) Apply(resp *fasthttp.Response) error {
	if time.Now().After(c.ExpireTS) {
		return ErrCacheItemExpired
	}

	_, err := resp.BodyWriter().Write(c.Body)
	if err != nil {
		return errors.Wrap(err)
	}

	resp.Header.Set(`Content-Type`, c.ContentType)
	resp.Header.Set(`X-Cached-Response`, `true`)
	resp.SetStatusCode(c.StatusCode)
	return nil
}
