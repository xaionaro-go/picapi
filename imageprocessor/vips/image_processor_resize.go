package imageprocessorimagick

import (
	"io"

	"github.com/davidbyttow/govips/pkg/vips"

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

	_, vipsImageType, err := vips.NewTransform().
		Load(in).
		Resize(int(toWidth), int(toHeight)).
		Output(out).
		Apply()
	if err != nil {
		return
	}

	switch vipsImageType {
	case vips.ImageTypeJPEG:
		imageFormat = imageprocessorcommon.ImageFormatJPEG
	default:
		imageFormat = imageprocessorcommon.ImageFormatOther
	}

	return
}
