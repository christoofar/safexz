package main

import (
	"fmt"

	"github.com/christoofar/safexz"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {

	var pass uint64
	p := message.NewPrinter(language.English)

	println("Starting compression...")

	for {
		pass++
		safexz.CompressFileWithProgress("/home/christoofar/ISO/debian.iso", fmt.Sprintf("/home/christoofar/VMBackups/debian%v.xz", pass), func(readByteCount uint64, decodedByteCount uint64) {
			print(p.Sprintf("\rPass: %v Read bytes: %v \tCompressed bytes: %v", pass, readByteCount, decodedByteCount))
		}, safexz.CompressionFullPowerMax)
		println()
		println("Compression complete.")
		println("Decompressing...")
		safexz.DecompressFileWithProgress(fmt.Sprintf("/home/christoofar/VMBackups/debian%v.xz", pass), fmt.Sprintf("/home/christoofar/VMBackups/debian%v.iso", pass), func(readByteCount uint64, decodedByteCount uint64) {
			print(p.Sprintf("\rPass: %v Read bytes: %v \tDecompressed bytes: %v", pass, readByteCount, decodedByteCount))
		})

		if pass == 100 {
			break
		}
	}
}
