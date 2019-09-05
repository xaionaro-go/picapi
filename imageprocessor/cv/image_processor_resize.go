package imageprocessorimagick

import (
	"image"
	"io"
	"io/ioutil"

	"gocv.io/x/gocv"

	"github.com/xaionaro-go/errors"
	imageprocessorcommon "github.com/xaionaro-go/picapi/imageprocessor/common"
)

var (
	ErrMatIsEmpty = errors.New(`mat is empty`)
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

	src, err := gocv.IMDecode(blob, -1)
	if err != nil {
		return
	}
	if src.Empty() {
		err = ErrMatIsEmpty
		return
	}

	dst := gocv.NewMat()
	gocv.Resize(src, &dst, image.Point{
		X: int(toWidth),
		Y: int(toHeight),
	}, 0, 0, gocv.InterpolationLanczos4)

	result, err := gocv.IMEncode(gocv.JPEGFileExt, dst)
	if err != nil {
		return
	}
	imageFormat = imageprocessorcommon.ImageFormatJPEG

	if dst.Empty() {
		err = ErrMatIsEmpty
		return
	}

	_, err = out.Write(result)
	if err != nil {
		return
	}

	src.Close()
	dst.Close()

	return
}
