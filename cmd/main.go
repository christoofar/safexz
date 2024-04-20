package main

import (
	"fmt"

	"github.com/christoofar/safexz"
)

func main() {

	// var pass uint64
	// p := message.NewPrinter(language.English)

	println("Starting compression...")

	compressedData, err := safexz.CompressBytes([]byte("Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!"))
	if err != nil {
		fmt.Println("Error compressing string:", err)
		return
	}
	fmt.Println("Compressed string 'Hello, World!':", compressedData)

	compressedString, err := safexz.CompressString("Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!")
	if err != nil {
		fmt.Println("Error compressing string:", err)
		return
	}
	fmt.Println("Compressed string 'Hello, World!':", []byte(compressedString))

	decompressedBytes, err := safexz.DecompressBytes(compressedData)
	if err != nil {
		fmt.Println("Error decompressing bytes:", err)
		return
	}
	fmt.Println("Decompressed bytes into string for hello world:", string(decompressedBytes))

	// for {
	// 	pass++
	// 	safexz.CompressFileWithProgress("/home/christoofar/ISO/debian.iso", fmt.Sprintf("/home/christoofar/VMBackups/debian%v.xz", pass), func(readByteCount uint64, decodedByteCount uint64) {
	// 		print(p.Sprintf("\rPass: %v Read bytes: %v \tCompressed bytes: %v", pass, readByteCount, decodedByteCount))
	// 	}, safexz.CompressionFullPowerMax)
	// 	println()
	// 	println("Compression complete.")
	// 	println("Decompressing...")
	// 	safexz.DecompressFileWithProgress(fmt.Sprintf("/home/christoofar/VMBackups/debian%v.xz", pass), fmt.Sprintf("/home/christoofar/VMBackups/debian%v.iso", pass), func(readByteCount uint64, decodedByteCount uint64) {
	// 		print(p.Sprintf("\rPass: %v Read bytes: %v \tDecompressed bytes: %v", pass, readByteCount, decodedByteCount))
	// 	})
	// 	println()
	// 	println("Decompression complete.")

	// 	if pass == 100 {
	// 		break
	// 	}
	// }
}
