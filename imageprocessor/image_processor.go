package imageprocessor

import (
	imageprocessorcommon "github.com/xaionaro-go/picapi/imageprocessor/common"
	imageprocessorgraphicsmagick "github.com/xaionaro-go/picapi/imageprocessor/imaging"
)

// ImageProcessor is a tool to manipulate images
type ImageProcessor = imageprocessorcommon.ImageProcessor

// NewImageProcessor creates an instance of `ImageProcessor` (image manipulation tool).
//
// Currently is uses a GraphicsMagick-based implementation (`github.com/xaionaro-go/picapi/imageprocessor/graphicsmagick`).
func NewImageProcessor() ImageProcessor {
	return imageprocessorgraphicsmagick.NewImageProcessor()
}
