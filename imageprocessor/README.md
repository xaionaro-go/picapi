```
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/picapi/imageprocessor/graphicsmagick
BenchmarkResize_smallPictures-8    	    3000	    441406 ns/op	   70789 B/op	      31 allocs/op
BenchmarkResize_mediumPictures-8   	     100	  18820133 ns/op	 4427312 B/op	      33 allocs/op
PASS
ok  	github.com/xaionaro-go/picapi/imageprocessor/graphicsmagick	4.441s
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/picapi/imageprocessor/imagick
BenchmarkResize_smallPictures-8    	    2000	    601050 ns/op	   75648 B/op	      31 allocs/op
BenchmarkResize_mediumPictures-8   	      50	  26750531 ns/op	 4933306 B/op	      33 allocs/op
PASS
ok  	github.com/xaionaro-go/picapi/imageprocessor/imagick	3.994s
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/picapi/imageprocessor/imaging
BenchmarkResize_smallPictures-8    	    2000	    700154 ns/op	  207362 B/op	     136 allocs/op
BenchmarkResize_mediumPictures-8   	      20	  58144684 ns/op	 8597955 B/op	      98 allocs/op
PASS
ok  	github.com/xaionaro-go/picapi/imageprocessor/imaging	4.190s
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/picapi/imageprocessor/nfntresize
BenchmarkResize_smallPictures-8    	    2000	    633334 ns/op	  182146 B/op	     123 allocs/op
BenchmarkResize_mediumPictures-8   	      20	  54274062 ns/op	11204054 B/op	      86 allocs/op
PASS
ok  	github.com/xaionaro-go/picapi/imageprocessor/nfntresize	3.836s
goos: linux
goarch: amd64
pkg: github.com/xaionaro-go/picapi/imageprocessor/vips
BenchmarkResize_smallPictures-8    	     200	   6544699 ns/op	   75591 B/op	      58 allocs/op
BenchmarkResize_mediumPictures-8   	     100	  19339423 ns/op	 4625103 B/op	      51 allocs/op
PASS
ok  	github.com/xaionaro-go/picapi/imageprocessor/vips	4.974s
```