package imageprocessorcommon

import "io"

// ImageProcessor is a tool to manipulate images
type ImageProcessor interface {
	// Resize reads an image from `in` and writes a resized image to `out`
	// (new image size: `toWidth`x`toHeight`).
	//
	// It returns image format and (if occured:) an error.
	Resize(in io.Reader, out io.Writer, toWidth, toHeight uint) (ImageFormat, error)
}
