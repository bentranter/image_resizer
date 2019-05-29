package main

import (
	"io/ioutil"
	"testing"
)

var img = read()

// read loads the file to test from the filesystem into memory.
//
// If you want to run this test yourself, change the filepath. Hopefully you
// don't have my avatar image on your computer at that path...
func read() []byte {
	file, err := ioutil.ReadFile("/Users/ben/Documents/ben-avatar-2018.jpeg")
	if err != nil {
		panic("failed to read file: " + err.Error())
	}
	return file
}

func BenchmarkResize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		resize(img)
	}
}
