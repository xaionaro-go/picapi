package httpserver

import (
	"context"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xaionaro-go/picapi/httpserver/helpers"
	imageprocessorcommon "github.com/xaionaro-go/picapi/imageprocessor/common"
)

type dummyImageProcessor struct{}

func (proc *dummyImageProcessor) Resize(
	in io.Reader,
	out io.Writer,
	toWidth, toHeight uint,
) (imageFormat imageprocessorcommon.ImageFormat, err error) {
	imageFormat = imageprocessorcommon.ImageFormatJPEG
	return
}

func TestHTTPServerStartStop(t *testing.T) {
	logger := log.New(os.Stderr, ``, 0)

	srv, err := NewHTTPServer(&dummyImageProcessor{}, logger, logger)
	assert.NoError(t, err)

	err = srv.Start(context.Background(), "")
	assert.NoError(t, err)

	err = srv.Stop()
	assert.NoError(t, err)

	err = srv.Start(context.Background(), "")
	assert.NoError(t, err)

	err = srv.Stop()
	assert.NoError(t, err)

	err = srv.Start(context.Background(), "")
	assert.NoError(t, err)

	func() {
		defer helpers.RecoverPanic(&err)
		err = srv.Start(context.Background(), "")
	}()
	assert.Error(t, err)
}
