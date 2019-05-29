# Image Resizer (and Optimizer)

A server that optimizes and resizes JPEGs (nothing else).

### Usage

After cloning this repo,

```
go run main.go
```

If this fails, ensure that you have cgo enabled, and your machine is supported by [lilliput](https://github.com/discordapp/lilliput).

### Benchmarks

A 1MB image is resized to 50% of it's original size while preserving its aspect ratio, and has a compression factor of 5x (from 1MB to 200KB). This resizing operation take on average 340ms, according to the benchmark below.

```
goos: darwin
goarch: amd64
pkg: github.com/bentranter/image_resizer
BenchmarkResize-4   	      10	 339673906 ns/op	542114793 B/op	      27 allocs/op
PASS
ok  	github.com/bentranter/image_resizer	3.853s
```
