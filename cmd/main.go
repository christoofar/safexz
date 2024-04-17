package main

import (
	"fmt"

	"github.com/christoofar/safexz"
)

func main() {

	// safexz.CompressFileWithProgress("../test/canterbury-corpus/large/world192.txt", func(decodedByteCount uint64) {
	// 	print("\rCompressed bytes:", decodedByteCount)
	// })

	safexz.CompressFileWithProgress("debian.iso", "debian.xz", func(readByteCount uint64, decodedByteCount uint64) {
		print(fmt.Sprintf("\rRead bytes: %v \tCompressed bytes: %v", readByteCount, decodedByteCount))
	})
}
