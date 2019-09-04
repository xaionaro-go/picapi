package imageprocessorimagick

import (
	"image/jpeg"
	"io"

	"github.com/nfnt/resize"

	"github.com/xaionaro-go/errors"
	imageprocessorcommon "github.com/xaionaro-go/picapi/imageprocessor/common"
)

// Resize reads an image from `in` and writes a resized image to `out`
// (new image size: `toWidth`x`toHeight`).
//
// It returns image format and (if occured:) an error.
func (proc *ImageProcessor) Resize(
	in io.Reader,
	out io.Writer,
	toWidth, toHeight uint,
) (imageFormat imageprocessorcommon.ImageFormat, err error) {
	defer func() { err = errors.Wrap(err) }()

	if toWidth <= 0 ||
		toHeight <= 0 {
		err = imageprocessorcommon.ErrInvalidSize
		return
	}

	image, err := jpeg.Decode(in)
	if err != nil {
		return
	}

	image = resize.Resize(toWidth, toHeight, image, resize.Lanczos2)
	imageFormat = imageprocessorcommon.ImageFormatJPEG

	err = jpeg.Encode(out, image, nil)
	if err != nil {
		return
	}

	return
}
