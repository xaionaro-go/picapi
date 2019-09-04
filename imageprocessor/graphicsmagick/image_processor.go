package imageprocessorgraphicsmagick

import (
	"github.com/gographics/gmagick"
)

func init() {
	gmagick.Initialize()
}

// ImageProcessor is a tool to manipulate images
type ImageProcessor struct {
}

// NewImageProcessor creates an instance of `ImageProcessor` (image manipulation tool).
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}
