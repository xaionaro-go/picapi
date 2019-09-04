[![Coverage Status](https://coveralls.io/repos/github/xaionaro-go/picapi/badge.svg?branch=master)](https://coveralls.io/github/xaionaro-go/picapi?branch=master)

# Introduction

"PicAPI" is an HTTP-server which implements a method to download a JPEG image and resize it.

# API

### `GET /resize`

| GET-parameter | description | required? |
| -------------:|:----------- | --------- |
| url           | An URL to download the image | Y |
| width         | resulting width (in pixels) | Y |
| height        | resulting height (in pixels) | Y |

#### Response

A JPEG image will be returned as the body.

If the response is returned from a cache then a header `X-Cached-Response` will also be returned.

#### cURL example:
```sh
curl -s --output /tmp/resized.jpg 'http://localhost:8486/resize?width=1230&height=200&url=https%3A%2F%2Fupload.wikimedia.org%2Fwikipedia%2Fcommons%2Fc%2Fc7%2FJPS-sample.jpg'
```

# Quick start

### Install build dependencies (Debian/Ubuntu)

```sh
sudo apt-get install libgraphicsmagick1-dev
```

If you're unable to use external libraries (like `graphicsmagick`) then you can switch to a native implementation (`imageprocessorimaging` or `imageprocessornfntresize` instead of `imageprocessorgraphicsmagick`) in file `imageprocessor/image_processor.go`.
