package safexz

import (
	"bytes"
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

// Construct some arbitrary string and compress it.
func TestXZReaderArbitraryString(t *testing.T) {
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
