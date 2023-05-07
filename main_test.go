package main

import "testing"

func BenchmarkListFiles(b *testing.B) {
	e := newFileList()
	e.listFiles("/Users/will/Library", false, "xxx.txt")
}
