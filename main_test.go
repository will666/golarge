package main

import "testing"

func BenchmarkListFiles(b *testing.B) {
	e := newFileList()
	e.listFiles("/", false, "xxx.txt", false, false)
}
