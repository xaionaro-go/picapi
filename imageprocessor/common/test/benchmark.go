package test

import (
	"bytes"
	"sync"
	"testing"

	imageprocessorcommon "github.com/xaionaro-go/picapi/imageprocessor/common"
)

// CommonBenchmarkResize is a function to run a benchmark for method `Resize`
// of an abstract ImageProcessor implementation.
func CommonBenchmarkResize(
	proc imageprocessorcommon.ImageProcessor,
	b *testing.B,
	widthMin, widthMax, heightMin, heightMax uint,
) {
	var items []Item
	for _, item := range Dataset() {
		if item.InputWidth < widthMin || item.InputWidth > widthMax ||
			item.InputHeight < heightMin || item.InputHeight > heightMax ||
			item.OutputWidth < widthMin || item.OutputWidth > widthMax ||
			item.OutputHeight < heightMin || item.OutputWidth > heightMax {
			continue
		}

		items = append(items, item)
	}

	bufPull := sync.Pool{New: func() interface{} {
		return &bytes.Buffer{}
	}}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, item := range items {
				reader := bytes.NewReader(item.InputFile)
				writer := bufPull.Get().(*bytes.Buffer)
				_, err := proc.Resize(reader, writer, item.OutputWidth, item.OutputHeight)
				if err != nil {
					b.Error(err)
				}
				writer.Reset()
				bufPull.Put(writer)
			}
		}
	})
}

// CommonBenchmarkResize_smallPictures is a function to run a benchmark for method `Resize`
// of an abstract ImageProcessor implementation on small images.
func CommonBenchmarkResize_smallPictures(proc imageprocessorcommon.ImageProcessor, b *testing.B) {
	CommonBenchmarkResize(proc, b, 1, 100, 1, 100)
}

// CommonBenchmarkResize_mediumPictures is a function to run a benchmark for method `Resize`
// of an abstract ImageProcessor implementation on medium images.
func CommonBenchmarkResize_mediumPictures(proc imageprocessorcommon.ImageProcessor, b *testing.B) {
	CommonBenchmarkResize(proc, b, 100, 1000, 100, 1000)
}

// CommonBenchmarkResize_largePictures is a function to run a benchmark for method `Resize`
// of an abstract ImageProcessor implementation on large images.
func CommonBenchmarkResize_largePictures(proc imageprocessorcommon.ImageProcessor, b *testing.B) {
	CommonBenchmarkResize(proc, b, 1000, 10000, 1000, 10000)
}
