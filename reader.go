package safexz

import (
	"io"

	internal "github.com/christoofar/safexz/internal/common"
)

type XZReader struct {
	io.Reader
	inputchan  chan []byte
	outputchan chan []byte
	started    bool
}

func (r *XZReader) Read(p []byte) (n int, err error) {
	if !r.started {
		// Start moving the reader data into the decompressor
		go func() {
			for {
				data := make([]byte, 1024)
				n, err := r.Reader.Read(data)
				if err != nil {
					close(r.inputchan)
					return
				}
				r.inputchan <- data[:n]
			}
		}()
		r.started = true
		// Start the decompressor
		internal.DecompressIn(r.inputchan, r.outputchan)
	}
	// Get the data from the decompressor
	data, ok := <-r.outputchan
	if !ok {
		return 0, io.EOF
	}
	n = copy(p, data)
	return n, nil
}

// Close closes the reader and the underlying channels.
func (r *XZReader) Close() error {
	// The underying channels close themselves and the memory grabbed by C is freed internally,
	// so we really don't need to do anything here.   It's just good practicy to have a close method.
	return nil
}

// NewReader creates a new XZReader that reads from r.   The data represented by r should have been compressed with XZ or LZMA.
func NewReader(r io.Reader) *XZReader {
	return &XZReader{Reader: r, inputchan: make(chan []byte), outputchan: make(chan []byte)}
}
