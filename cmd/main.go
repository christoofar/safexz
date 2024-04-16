package main

import "github.com/christoofar/safexz"

func main() {

	safexz.CompressFileWithProgress("../test/canterbury-corpus/large/world192.txt", func(progress float64) {})

}
