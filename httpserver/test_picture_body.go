package httpserver

import (
	"bytes"
	"encoding/base64"
	"io"
	"strings"
)

var (
	testPictureBodyValue = []byte{}
)

func init() {
	testPictureBodyValue = make([]byte, 1024)
	n, err := base64.StdEncoding.Decode(testPictureBodyValue, []byte(strings.TrimSpace(``+
		`/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAMCAgICAgMCAgIDAwMDBAYEBAQEBAgGBgUGCQgKCgkI`+
		`CQkKDA8MCgsOCwkJDRENDg8QEBEQCgwSExIQEw8QEBD/2wBDAQMDAwQDBAgEBAgQCwkLEBAQEBAQ`+
		`EBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBD/wAARCAAKAAoDAREA`+
		`AhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAABQf/xAAlEAABAwMDBAMBAAAAAAAAAAADAQIEBQYR`+
		`EhMUAAchIhYxQUL/xAAYAQADAQEAAAAAAAAAAAAAAAAEBQYCA//EACkRAAIBAwMCBAcAAAAAAAAA`+
		`AAECAwQFEQASIQYTFCIjMSQyM0FRYWL/2gAMAwEAAhEDEQA/ACaPLi23Ippi0Gzo0aUOLGWp8IEZ`+
		`BGdPRJQWhcoxOKwxzDaozNEwbTbZiaWPZ3ttCaq00lNVxzu5ZWQ5OFk2blwXlLqhT1dixJM0qiQR`+
		`RK+5t09RWVVA/hVDx7o5t21VcxbPSk3KIhF9LL9x3DyPFHGzE9t2ob7A4gOF3uuOFH227MblPDss`+
		`x6s29CaNKYTTjxjH50jmrbTDI0dRZ7UzgkMxiqSSR7klYmUknklWKn7EjnRk/XUdNK0C9N084Ukd`+
		`xuwrSY43sDTEhn+YgngnGn+6QhUOHZpqKNtPJKvQrDviptKVvylBYcrcak20RmF/n1+vHUR0RX1d`+
		`TdauSaVmaFUSMliSiGA+RCT5U/kYH60Ld/hulZq2DyyxGYI44ZAQ2QjDlQcnIBHufzqAXV3Dv+k3`+
		`RWKVSr4uCFChT5EeNGj1MwxBEwjmsYxjXIjWtRERERMIiIidWF6u1wp7lUQwzuqq7gAMwAAYgAAH`+
		`AAHsNMLp05Zo66ZEpIgA7AARpgDJ4HGv/9k=`)))
	if err != nil {
		panic(err)
	}
	testPictureBodyValue = testPictureBodyValue[:n]
}

type testPictureBodyStruct struct {
	*bytes.Reader
}

func testPictureBody() io.ReadCloser {
	return &testPictureBodyStruct{bytes.NewReader(testPictureBodyValue)}
}

func (s *testPictureBodyStruct) Close() (err error) {
	return
}
