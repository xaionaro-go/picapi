package imageprocessorimagick

import (
	"gopkg.in/gographics/imagick.v2/imagick"
)

func init() {
	imagick.Initialize()
}

// ImageProcessor is a tool to manipulate images
type ImageProcessor struct {
}

// NewImageProcessor creates an instance of `ImageProcessor` (image manipulation tool).
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}
