package imageprocessorgraphicsmagick

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gographics/gmagick"
)

func init() {
	gmagick.Initialize()

	signal.Ignore(os.Interrupt)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		<-ch
		os.Exit(int(syscall.EINTR))
	}()
}

// ImageProcessor is a tool to manipulate images
type ImageProcessor struct {
}

// NewImageProcessor creates an instance of `ImageProcessor` (image manipulation tool).
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}
