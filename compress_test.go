package safexz

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"
)

// TestCompressString tests the CompressString function
func TestCompressString(t *testing.T) {
	compressedString, err := CompressString("Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!")
	if err != nil {
		t.Errorf("Error compressing string: %v", err)
	}
	if len(compressedString) == 0 {
		t.Errorf("Compressed string is empty")
	}
}

// TestCompressBytes tests the CompressBytes function
func TestCompressBytes(t *testing.T) {
	compressedData, err := CompressBytes([]byte("Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!"))
	if err != nil {
		t.Errorf("Error compressing bytes: %v", err)
	}
	if len(compressedData) == 0 {
		t.Errorf("Compressed data is empty")
	}
}

// TestCompressFile tests the CompressFile function
func TestCompressFile(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}
	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz")
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")

}

// TestCompressFileSimple tests the CompressFile function with the CompressionSimple strategy.
func TestCompressFileSimple(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionSimple)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

// TestCompressFileSimpleBetter tests the CompressFile function with the CompressionSimpleBetter strategy.
func TestCompressFileSimpleBetter(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionSimpleBetter)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

func TestCompressFileSimpleMax(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionSimpleMax)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

// TestCompressFileSimpleFast tests the CompressFile function with the CompressionSimpleFast strategy.
func TestCompressFileSimpleFast(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionSimpleFast)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

func TestCompressFileMulti(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionMulti)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

func TestCompressFileMultiBetter(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionMultiBetter)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

// TestCompressFileSimpleFast tests the CompressFile function with the CompressionSimpleFast strategy.
func TestCompressFileMuliFast(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionMultiFast)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

func TestCompressFileMultiMax(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionMultiMax)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

// TestCompressFileFullPower tests the CompressFile function with the CompressionFullPower strategy.
func TestCompressFileFullPower(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionFullPower)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

// TestCompressFileFullPowerBetter tests the CompressFile function with the CompressionFullPowerBetter strategy.
func TestCompressFileFullPowerBetter(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionFullPowerBetter)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

// TestCompressFileFullPowerFast tests the CompressFile function with the CompressionFullPowerFast strategy.
func TestCompressFileFullPowerFast(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionFullPowerFast)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

// TestCompressFileFullPowerMax tests the CompressFile function with the CompressionFullPowerMax strategy.
func TestCompressFileFullPowerMax(t *testing.T) {
	// Create a test file
	f, err := os.Create("test.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// Write 1000 "Hello, World!" strings to the file
	for i := 0; i < 1000; i++ {
		f.WriteString("Hello, World! ")
	}
	f.Close()

	// Compress the file
	err = CompressFile("test.txt", "test.txt.xz", CompressionFullPowerMax)
	if err != nil {
		t.Errorf("Error compressing file: %v", err)
	}

	// The compressed file should be 140 bytes
	fi, err := os.Stat("test.txt.xz")
	if err != nil {
		t.Errorf("Error getting compressed file info: %v", err)
	}
	if fi.Size() != 140 {
		t.Errorf("Compressed file size is not 140 bytes")
	}

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
}

// TestCompressChain creates a scratch file using /test/large text files, takes a checksum then runs a chain of compressions
// and decompressions of the file, making sure none of the bytes moves or gets destroyed.
func TestCompressChain(t *testing.T) {
	// get the scratch files
	f1, err := os.Open("test/canterbury-corpus/large/bible.txt")
	if err != nil {
		t.Errorf("Error opening test file: %v", err)
	}
	f2, err := os.Open("test/canterbury-corpus/large/E.coli")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}
	f3, err := os.Open("test/canterbury-corpus/large/world192.txt")
	if err != nil {
		t.Errorf("Error creating test file: %v", err)
	}

	// generate a random number from 1-3
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(3) + 1

	// build a scratch file
	t.Log("Building scratch file from canterbury-corpus/large files")
	scratch, err := os.Create("scratch.txt")
	if err != nil {
		t.Errorf("Error creating scratch file: %v", err)
	}
	for i := 0; i < 50; i++ {
		switch randomNumber {
		case 1:
			f1.Seek(0, 0)
			f1.WriteTo(scratch)
		case 2:
			f2.Seek(0, 0)
			f2.WriteTo(scratch)
		case 3:
			f3.Seek(0, 0)
			f3.WriteTo(scratch)
		}
		randomNumber = rand.Intn(3) + 1
	}

	// close the files, start the compressions
	f1.Close()
	f2.Close()
	f3.Close()
	scratch.Close()

	// get the checksum of the scratch file
	checksum, err := ChecksumFile("scratch.txt")
	if err != nil {
		t.Errorf("Error getting checksum of scratch file: %v", err)
	}
	t.Logf("MD5 checksum of scratch file: %v", checksum)

	// do 5 compressions and decompressions
	for i := 0; i < 5; i++ {
		t.Logf("Starting compression %v", i)
		err = CompressFile("scratch.txt", "scratch.txt.xz", CompressionFullPowerFast)
		if err != nil {
			t.Errorf("Error compressing file: %v", err)
		}
		fi, _ := os.Stat("scratch.txt.xz")
		t.Logf("Compressed file size: %v", fi.Size())

		t.Logf("Starting decompression %v", i)
		err = DecompressFile("scratch.txt.xz", "scratch.txt")
		if err != nil {
			t.Errorf("Error decompressing file: %v", err)
		}

		// get the checksum of the scratch file
		checksum2, err := ChecksumFile("scratch.txt")
		if err != nil {
			t.Errorf("Error getting checksum of scratch file: %v", err)
		}
		t.Logf("MD5 checksum of decompressed scratch file: %v", checksum2)

		// compare the checksums
		if checksum != checksum2 {
			t.Errorf("Checksums do not match after compression and decompression")
		}

		if i != 4 {
			t.Log("Compressing again.")
		}
	}

	os.Remove("scratch.txt")
	os.Remove("scratch.txt.xz")

}

// ChecksumFile returns the MD5 checksum of a file
func ChecksumFile(s string) (string, error) {
	file, err := os.Open(s)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}
