package main

import (
	"fmt"

	"github.com/christoofar/safexz"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {

	// safexz.CompressFileWithProgress("../test/canterbury-corpus/large/world192.txt", func(decodedByteCount uint64) {
	// 	print("\rCompressed bytes:", decodedByteCount)
	// })

	// safexz.CompressFileWithProgress("debian.iso", "debian.xz", func(readByteCount uint64, decodedByteCount uint64) {
	// 	print(fmt.Sprintf("\rRead bytes: %v \tCompressed bytes: %v", readByteCount, decodedByteCount))
	// })

	var pass uint64
	p := message.NewPrinter(language.English)

	for {
		pass++
		safexz.CompressFileWithProgress("/home/christoofar/ISO/debian.iso", fmt.Sprintf("/home/christoofar/VMBackups/debian%v.xz", pass), func(readByteCount uint64, decodedByteCount uint64) {
			print(p.Sprintf("\rPass: %v Read bytes: %v \tCompressed bytes: %v", pass, readByteCount, decodedByteCount))
		}, safexz.CompressionMulti)
		println()
		return
	}
}
