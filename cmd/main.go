package main

import (
	"bytes"
	"fmt"
	"os"

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

	// check the files
	comp1, _ := os.Open("/home/christoofar/ISO/debian.iso")
	comp2, _ := os.Open("/home/christoofar/VMBackups/debian1")

	count := 0
	for {
		count += 1024
		comp1Bytes := make([]byte, 1024)
		comp2Bytes := make([]byte, 1024)

		_, err1 := comp1.Read(comp1Bytes)
		_, err2 := comp2.Read(comp2Bytes)

		if err1 != nil || err2 != nil {
			break
		}

		if !bytes.Equal(comp1Bytes, comp2Bytes) {
			println("Files are not equal", count)
			os.WriteFile("/home/christoofar/VMBackups/comp1.bin", comp1Bytes, 0644)
			os.WriteFile("/home/christoofar/VMBackups/comp2.bin", comp2Bytes, 0644)
			break
		}
	}

	println("Starting compression...")

	for {
		pass++
		safexz.CompressFileWithProgress("/home/christoofar/ISO/debian.iso", fmt.Sprintf("/home/christoofar/VMBackups/debian%v.xz", pass), func(readByteCount uint64, decodedByteCount uint64) {
			print(p.Sprintf("\rPass: %v Read bytes: %v \tCompressed bytes: %v", pass, readByteCount, decodedByteCount))
		}, safexz.CompressionFullPowerFast)
		println()
		return
	}
}
