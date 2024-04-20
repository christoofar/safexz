package safexz

import (
	"os"
	"testing"
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

// TestCompressStringSimple tests the CompressString function
func TestCompressBytes(t *testing.T) {
	compressedData, err := CompressBytes([]byte("Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!Hello, World!"))
	if err != nil {
		t.Errorf("Error compressing bytes: %v", err)
	}
	if len(compressedData) == 0 {
		t.Errorf("Compressed data is empty")
	}
}

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
