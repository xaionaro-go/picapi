package test

import (
	"bytes"
	"image/jpeg"
	"testing"

	"github.com/stretchr/testify/assert"

	imageprocessorcommon "github.com/xaionaro-go/picapi/imageprocessor/common"
)

// CommonTestResize verifies if the method `Resize` of an ImageProcessor works correcty.
func CommonTestResize(proc imageprocessorcommon.ImageProcessor, t *testing.T) {
	dataset := Dataset()

	// Positive testing:

	for _, item := range dataset {
		if item.InputWidth > 2000 {
			// A have not enough RAM on my laptop to use test high resolutions :(
			continue
		}
		if item.OutputWidth < item.InputWidth || item.OutputHeight < item.InputHeight {
			// images-diff will be high on scaling-down a noise
			continue
		}

		reader := bytes.NewReader(item.InputFile)
		writer := &bytes.Buffer{}

		imageFormat, err := proc.Resize(reader, writer, item.OutputWidth, item.OutputHeight)
		if err != nil {
			t.Error(err)
		}
		if imageFormat != imageprocessorcommon.ImageFormatJPEG {
			t.Errorf(`unexpected image format: %v`, imageFormat)
		}

		diff := diffImages(item.InputFile, writer.Bytes())
		if diff == 0 {
			t.Errorf(`no difference between images`)
		}
		if diff > 0.2 { // We're scaling noises, so huge diffs (up-to 0.2) is OK.
			t.Errorf(`images-diff is too high: %v`, diff)
		}

		width, height := getSize(writer.Bytes())
		if width != item.OutputWidth {
			t.Errorf(`wrong resulting width: %v (expected %v)`, width, item.OutputWidth)
		}

		if height != item.OutputHeight {
			t.Errorf(`wrong resulting height: %v (expected %v)`, height, item.OutputHeight)
		}
	}

	// Negative testing:

	// no input
	//_, err := proc.Resize(nil, &bytes.Buffer{}, 100, 100)
	//assert.Error(t, err)
	// it's a valid panic-case, so it's commented-out

	// no output
	//_, err = proc.Resize(bytes.NewReader(dataset[0].InputFile), nil, 100, 100)
	//assert.Error(t, err)
	// it's a valid panic-case, so it's commented-out

	// invalid image
	_, err := proc.Resize(bytes.NewReader(make([]byte, 100)), &bytes.Buffer{}, 100, 100)
	assert.Error(t, err)

	// invalid size (the last case is commented-out because considered valid for a while)
	_, err = proc.Resize(bytes.NewReader(dataset[0].InputFile), &bytes.Buffer{}, 0, 100)
	assert.Error(t, err)
	_, err = proc.Resize(bytes.NewReader(dataset[0].InputFile), &bytes.Buffer{}, 100, 0)
	assert.Error(t, err)
	//_, err = proc.Resize(bytes.NewReader(dataset[0].InputFile), &bytes.Buffer{}, 100000, 100000)
	//assert.Error(t, err)
}

func min(args ...int) int {
	r := args[0]
	for _, v := range args[1:] {
		if v < r {
			r = v
		}
	}
	return r
}

func abs(a int) uint {
	if a < 0 {
		return uint(-a)
	}
	return uint(a)
}

func diffImages(aBytes, bBytes []byte) float64 {
	a, err := jpeg.Decode(bytes.NewReader(aBytes))
	panicIfError(err)

	b, err := jpeg.Decode(bytes.NewReader(bBytes))
	panicIfError(err)

	aSize := a.Bounds().Size()
	bSize := b.Bounds().Size()

	var totalDiff float64
	var count uint64

	minX := min(aSize.X, bSize.X)
	minY := min(aSize.Y, bSize.Y)
	for x := 0; x < minX; x++ {
		for y := 0; y < minY; y++ {
			aX := int(float64(aSize.X) * float64(x) / float64(minX))
			aY := int(float64(aSize.Y) * float64(y) / float64(minY))
			bX := int(float64(bSize.X) * float64(x) / float64(minX))
			bY := int(float64(bSize.Y) * float64(y) / float64(minY))

			aR, aG, aB, aA := a.At(aX, aY).RGBA()
			bR, bG, bB, bA := b.At(bX, bY).RGBA()

			totalDiff += float64(abs(int(aR)-int(bR))) / (1 << 16 /* maximal value */)
			totalDiff += float64(abs(int(aG)-int(bG))) / (1 << 16 /* maximal value */)
			totalDiff += float64(abs(int(aB)-int(bB))) / (1 << 16 /* maximal value */)
			totalDiff += float64(abs(int(aA)-int(bA))) / (1 << 16 /* maximal value */)

			count += 4 /* there was 4 diffs above ^ */
		}
	}

	return float64(totalDiff) / float64(count)
}

func getSize(b []byte) (width, height uint) {
	img, err := jpeg.Decode(bytes.NewReader(b))
	panicIfError(err)

	size := img.Bounds().Size()
	return uint(size.X), uint(size.Y)
}
