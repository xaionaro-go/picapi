package imageprocessor

import (
	imageprocessorcommon "github.com/xaionaro-go/picapi/imageprocessor/common"
)

// ImageFormat is an enumeration data-type for image formats
type ImageFormat = imageprocessorcommon.ImageFormat

const (
	// ImageFormatUndefined is a reserved ImageFormat for cases when it's undefined
	// So if anybody will have an ImageFormat with zero-value (not defined, yet) it won't receive
	// a mis-leading value like ImageFormatJPEG.
	ImageFormatUndefined = imageprocessorcommon.ImageFormatUndefined

	// ImageFormatJPEG is an ImageFormat for JPEGs
	ImageFormatJPEG = imageprocessorcommon.ImageFormatJPEG

	// ImageFormatOther is an ImageFormat for all image formats not listed here in this enumeration
	ImageFormatOther = imageprocessorcommon.ImageFormatOther
)
