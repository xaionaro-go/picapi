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

### Install

```sh
go get github.com/xaionaro-go/picapi/cmd/picapid
```

### Run

```sh
PICAPI_LOGGING_LEVEL=debug $(go env GOPATH)/bin/picapid
```

# Performance test

```sh
$(go env GOPATH)/bin/picapid &
$(go env GOPATH)/bin/gobench -u 'http://localhost:8486/resize?width=10&height=10&url=https%3A%2F%2Fupload.wikimedia.org%2Fwikipedia%2Fcommons%2Fc%2Fc7%2FJPS-sample.jpg' -t 10
```

```
Dispatching 100 clients
Waiting for results...

Requests:                           856792 hits
Successful requests:                856792 hits
Network failed:                          0 hits
Bad requests failed (!2xx):              0 hits
Successful requests rate:            85679 hits/sec
Read throughput:                  78311428 bytes/sec
Write throughput:                 16966600 bytes/sec
Test time:                              10 sec
```

```
cpufreq-info  | grep 'current policy' | head -1; echo; grep 'model name' /proc/cpuinfo  | head -1

  current policy: frequency should be within 800 MHz and 3.70 GHz.
  
model name      : Intel(R) Core(TM) i7-4800MQ CPU @ 2.70GHz
```

# Packages

* `imageprocessor` is a wrapper around different image manipulation tools (like imagemagick). Currently we use `graphicsmagick` as the backend.
* `httpserver` is a wrapper around an http router (currently we use `fasthttprouter`) and `imageprocessor` to serve incoming HTTP requests.
* `main` (`cmd/picapid`) is the entry-point/executable package.
* `config` is just a structure to configure the `main`.
