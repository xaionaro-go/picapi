```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/picapi/imageprocessor/cv
BenchmarkResize_smallPictures-8    	   10000	    224289 ns/op	   76152 B/op	      19 allocs/op
BenchmarkResize_mediumPictures-8   	     100	  16256349 ns/op	 4821099 B/op	      25 allocs/op
PASS
ok  	github.com/xaionaro-go/picapi/imageprocessor/cv	5.055s
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/picapi/imageprocessor/graphicsmagick
BenchmarkResize_smallPictures-8    	    3000	    430436 ns/op	   70794 B/op	      31 allocs/op
BenchmarkResize_mediumPictures-8   	     100	  18352913 ns/op	 4429245 B/op	      33 allocs/op
PASS
ok  	github.com/xaionaro-go/picapi/imageprocessor/graphicsmagick	4.331s
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/picapi/imageprocessor/imagick
BenchmarkResize_smallPictures-8    	    2000	    609666 ns/op	   75639 B/op	      31 allocs/op
BenchmarkResize_mediumPictures-8   	      50	  27350344 ns/op	 4934413 B/op	      33 allocs/op
PASS
ok  	github.com/xaionaro-go/picapi/imageprocessor/imagick	4.089s
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/picapi/imageprocessor/imaging
BenchmarkResize_smallPictures-8    	    2000	    770229 ns/op	  207365 B/op	     137 allocs/op
BenchmarkResize_mediumPictures-8   	      20	  59357494 ns/op	 8593362 B/op	      98 allocs/op
PASS
ok  	github.com/xaionaro-go/picapi/imageprocessor/imaging	4.432s
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/picapi/imageprocessor/nfntresize
BenchmarkResize_smallPictures-8    	    2000	    653626 ns/op	  182103 B/op	     123 allocs/op
BenchmarkResize_mediumPictures-8   	      20	  61011215 ns/op	11213379 B/op	      87 allocs/op
PASS
ok  	github.com/xaionaro-go/picapi/imageprocessor/nfntresize	4.183s
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/picapi/imageprocessor/vips
BenchmarkResize_smallPictures-8    	     200	   6913212 ns/op	   75592 B/op	      58 allocs/op
BenchmarkResize_mediumPictures-8   	     100	  19660814 ns/op	 4625096 B/op	      51 allocs/op
PASS
ok  	github.com/xaionaro-go/picapi/imageprocessor/vips	5.547s
```