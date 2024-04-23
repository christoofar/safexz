package safexz

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXZWriter(t *testing.T) {
	// This is a simple test to make sure the writer doesn't panic.
	// The writer is tested in the internal package.
	w := NewWriter(nil)
	if w == nil {
		t.Error("NewWriter returned nil")
	}

	// w.Close() <- Guess what you can't do this because there's no way to test that the write is done on an empty writer,
	// as Close() will close the internal send channel so LZMA knows compression is ending.  You can read me bitching about
	// this in writer.go
}

// Construct some arbitrary string and compress it through the byte writer ABI.  The work is verified by decompressing the bytes
// without using XZReader.
func TestXZWriterArbitraryString(t *testing.T) {
	testString := "This is a test string that we will compress and then decompress."

	bytesToRead := bytes.NewReader([]byte(testString))
	compressedBytes := new(bytes.Buffer)
	xzwriter := NewWriter(compressedBytes)

	for {
		buffer := make([]byte, 5)  // For punishment, we'll set a ridiculously low buffer size to force cycling.
		n, err := bytesToRead.Read(buffer)
		if err != nil {
			xzwriter.Close()
			break
		}
		if n > 0 {
			xzwriter.Write(buffer[:n])
		}
	}

	t.Log("Compressed bytes:", compressedBytes.Bytes())

	// Now decompress the bytes
	decompressedBytes, err := DecompressBytes(compressedBytes.Bytes())
	if err != nil {
		t.Error("Error decompressing bytes:", err)
	}

	t.Log("Decompressed bytes:", decompressedBytes)
	assert.Equal(t, testString, string(decompressedBytes), "Decompressed string does not match original string.")

}

// Construct some arbitrary string and compress it without using the byte writer ABI.  Then read the compressed string back
// using the XZReader sitting on top of an io.Reader
func TestXZReaderArbitraryString(t *testing.T) {
	testString := "This is a test string that we will compress and then decompress."

	compressedBytes, err := CompressBytes([]byte(testString))
	if err != nil {
		t.Error("Error compressing string:", err)
	}

	t.Log("Compressed bytes:", compressedBytes)

	// Make a reader from the compressed bytes which will be fed to the XZReader for decompression
	compressedReader := bytes.NewReader(compressedBytes)
	// Pull the compressedReader through the XZReader
	decompressedReader := NewReader(compressedReader)

	// Read the decompressed bytes
	decompressedBytes := make([]byte, len(testString)) // Adjust the buffer size as needed
	n := 0
	for {
		n, err = decompressedReader.Read(decompressedBytes)
		if n == 0 {
			break
		}
	}
	if err != nil {
		t.Error("Error reading decompressed bytes:", err)
	}

	t.Log("Decompressed bytes:", decompressedBytes)
	assert.Equal(t, testString, string(decompressedBytes), "Decompressed string does not match original string.")
}

// Complex test case.
//
// Test the XZReader and XZWriter together.  This is a simple test to make sure the reader and writer work together.
// We'll start with loading a test file, compressing it to an .xz file, then read it up with the XZReader and compress it
// yet again with XZWriter, reassigned to the same in-place variable.
//
// The final compressed file should be the same as the original compressed file.
func TestXZReaderWriter(t *testing.T) {

	// Let's begin by compressing a copy of the Bible
	err := CompressFile("test/canterbury-corpus/large/bible.txt", "test.txt.xz", CompressionMulti)
	if err != nil {
		t.Error("Error compressing file:", err)
	}

	// Decompress the Bible for comparison later
	err = DecompressFile("test.txt.xz", "test.txt")
	if err != nil {
		t.Error("Error decompressing file:", err)
	}

	// Now let's test the XZReader... did it decompress the file the same way
	// as the DecompressFile function which uses the direct way?
	readmem := bytes.Buffer{}
	readfile, err := os.Open("test.txt.xz")
	if err != nil {
		t.Error("Error opening compressed file:", err)
	}
	readreader := NewReader(readfile)  // This is the XZReader
	_, err = io.Copy(&readmem, readreader) // Read the compressed file into a buffer
	readreader.Close()
	if err != nil {
		t.Error("Error reading compressed file:", err)
	}

	// Now that we have the XZReader's version in memory, push it back out to a file
	// with XZWriter and compare the two files.
	writefile, err := os.Create("test2.txt.xz")
	writer := NewWriter(writefile)  // This is the XZWriter
	if err != nil {
		t.Error("Error creating file for comparison:", err)
	}
	_, err = io.Copy(writer, &readmem)
	writer.Close() /* <--- YOU MUST REMEMBER TO DO THIS WITH THE XZWRITER!!!!!!!!!!!!!!!!!!!!!!!!!!! */
	if err != nil {
		t.Error("Error writing file for comparison:", err)
	}
	writefile.Close()

	// Now decompress text2.txt.xz into test2.txt
	err = DecompressFile("test2.txt.xz", "test2.txt")
	if err != nil {
		t.Error("Error decompressing file:", err)
	}

	// Are text.txt and test2.txt the same?
	original, _ := os.ReadFile("test.txt")
	compare, _ := os.ReadFile("test2.txt")
	assert.Equal(t, original, compare, "Decompressed files do not match.")


	// Now compress the file again.  This test case proves that the XZWriter var can be re-used.
	writefile, err = os.Create("test2.txt.xz")
	if err != nil {
		t.Error("Error creating compressed file:", err)
	}
	writer = NewWriter(writefile)  // This is the XZWriter
	_, err = io.Copy(writer, &readmem) // Write the compressed file to the new file
	writer.Close()  // <--- YOU MUST REMEMBER TO DO THIS WITH THE XZWRITER!!!!!!!!!!!!!!!!!!!!!!!!!!!
	if err != nil {
		t.Error("Error writing compressed file:", err)
	}


	// Now, decompress both products
	err = DecompressFile("test.txt.xz", "test.txt")
	if err != nil {
		t.Error("Error decompressing file:", err)
	}


	// Are these the same?
	original, _ = os.ReadFile("test.txt")
	compare, _ = os.ReadFile("test2.txt")
	assert.Equal(t, original, compare, "Re-decompressed files do not match.")

	// Clean up
	os.Remove("test.txt")
	os.Remove("test.txt.xz")
	os.Remove("test2.txt")
	os.Remove("test2.txt.xz")
}


func TestXZReaderWeirdCaseTinyBuffers(t *testing.T) {
	// Let's create a file with 15KB of test data, and create a read buffer of 1 byte
	// to test the XZReader's ability to handle tiny buffers.
originalmem, _ := os.ReadFile("test/canterbury-corpus/large/bible.txt")
	err := CompressFile("test/canterbury-corpus/large/bible.txt", "test.txt.xz", CompressionMulti)
	if err != nil {
		t.Error("Error compressing file:", err)
	}

	// Decompress the Bible for comparison
	mem := bytes.NewBuffer(nil)
	mem.Grow(len(originalmem))
	osfile, err := os.Open("test.txt.xz")
	if err != nil {
		t.Error("Error opening compressed file:", err)
	}
	reader := NewReader(osfile)
	for {
		buffer := make([]byte, 3)  // Ridiculously small buffer read buffer to exercise the holdover code in XZReader
		n, err := reader.Read(buffer)
		if err != nil {
			break
		}
		if n > 0 {
			mem.Write(buffer[:n])
		}
	}
	osfile.Close()
	println(len(originalmem))

	// Now decompress the file the direct way
	bible, err := DecompressFileToMemory("test.txt.xz")
	if err != nil {
		t.Error("Error decompressing file:", err)
	}

	// Does bible match the contents of mem?
	assert.Equal(t, bible, mem.Bytes(), "Decompressed files do not match.")
}