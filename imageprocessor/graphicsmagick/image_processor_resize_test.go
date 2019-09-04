package imageprocessorgraphicsmagick

import (
	"testing"

	"github.com/xaionaro-go/picapi/imageprocessor/common/test"
)

func BenchmarkResize_smallPictures(b *testing.B) {
	test.CommonBenchmarkResize_smallPictures(NewImageProcessor(), b)
}

func BenchmarkResize_mediumPictures(b *testing.B) {
	test.CommonBenchmarkResize_mediumPictures(NewImageProcessor(), b)
}

/*func BenchmarkResize_largePictures(b *testing.B) {
	test.CommonBenchmarkResize_largePictures(NewImageProcessor(), b)
}*/

func TestResize(t *testing.T) {
	test.CommonTestResize(NewImageProcessor(), t)
}
