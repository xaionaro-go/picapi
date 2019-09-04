package imageprocessorimagick

import (
	"io"

	"github.com/disintegration/imaging"

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

	image, err := imaging.Decode(in)
	if err != nil {
		return
	}

	image = imaging.Resize(image, int(toWidth), int(toHeight), imaging.Lanczos)

	err = imaging.Encode(out, image, imaging.JPEG)
	if err != nil {
		return
	}

	imageFormat = imageprocessorcommon.ImageFormatJPEG

	return
}
