package imageprocessorcommon

import "strings"

// ImageFormat is an enumeration data-type for image formats
type ImageFormat uint

const (
	// ImageFormatUndefined is a reserved ImageFormat for cases when it's undefined.
	// So if anybody will have an ImageFormat with zero-value (not defined, yet) it won't receive
	// a mis-leading value like ImageFormatJPEG.
	ImageFormatUndefined = ImageFormat(iota)

	// ImageFormatJPEG is an ImageFormat for JPEGs
	ImageFormatJPEG

	// ImageFormatOther is an ImageFormat for all image formats not listed here in this enumeration
	ImageFormatOther
)

// ParseImageFormat parses a string with an image format to ImageFormat.
//
// Expected values: "JPEG"
func ParseImageFormat(formatString string) ImageFormat {
	switch strings.ToUpper(formatString) {
	case ``:
		return ImageFormatUndefined
	case `JPEG`:
		return ImageFormatJPEG
	default:
		return ImageFormatOther
	}
}

// String just implements Stringer for nicer logging
func (ifmt ImageFormat) String() string {
	switch ifmt {
	case ImageFormatUndefined:
		return `<undefined>`
	case ImageFormatJPEG:
		return `JPEG`
	case ImageFormatOther:
		return `<other>`
	default:
		panic(`should not happened`)
	}
}
