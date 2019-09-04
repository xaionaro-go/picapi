package imageprocessorgraphicsmagick

import (
	"io"
	"io/ioutil"

	"github.com/gographics/gmagick"

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

	blob, err := ioutil.ReadAll(in)
	if err != nil {
		return
	}

	wand := gmagick.NewMagickWand()
	defer wand.Destroy()

	err = wand.ReadImageBlob(blob)
	if err != nil {
		return
	}

	imageFormat = imageprocessorcommon.ParseImageFormat(wand.GetImageFormat())

	err = wand.ResizeImage(toWidth, toHeight, gmagick.FILTER_GAUSSIAN, 1)
	if err != nil {
		return
	}

	_, err = out.Write(wand.WriteImageBlob())
	if err != nil {
		return
	}

	return
}
