package safexz

import (
	"io"

	internal "github.com/christoofar/safexz/internal"
)

// XZReader reads an LZMA1 or LZMA2 compressed stream from the supplied source and yields the compressed data into a byte slice.  When the end of
// the compressed stream is reached, Read will return io.EOF.
type XZReader struct {
	io.Reader
	inputchan  chan []byte
	outputchan chan []byte
	started    bool
	holdover  []byte
}

// Read reads an LZMA1 or LZMA2 compressed stream from the supplied soure and yields the compressed data into a byte slice.
func (r *XZReader) Read(p []byte) (n int, err error) {
	if !r.started {
		go func() {
			// Start the decompressor
			internal.DecompressIn(r.inputchan, r.outputchan)
		}()

		// Start moving the reader data into the decompressor
		go func() {
			for {
				data := make([]byte, len(p))
				n, err := r.Reader.Read(data)
				if err != nil {
					close(r.inputchan)
					return
				}
				r.inputchan <- data[:n]
			}
		}()

		r.started = true
	}

	// If there are holdover bytes from the last read, we need to put them in p first.
	if len(r.holdover) > 0 {
		n = copy(p, r.holdover)
		// Shrink the holdover slice by n bytes, from the front of the slice.
		if n > 0 {
			r.holdover = r.holdover[n:]
		}
		return n, nil
	}

	// Get data from the decompressor.  Read has to be called again to get the next block.
	data, ok := <-r.outputchan

	// If there is more data in data than we can fit in p, we need to hold it over for the next read.
	if len(data) > len(p) {
		r.holdover = data[len(p):]
		data = data[:len(p)]
	}

	// But if less data came out than the size of p, we should try to pull more data from the decompressor.
	if len(data) < len(p) {
		limit := len(p) - len(data)	// The amount of data we can still fit in p
		// Get more data from the decompressor
		data2, ok := <-r.outputchan
		if ok {
			// If there is more data in data2 than we can fit in p, we need to hold it over for the next read.
			if len(data2) > limit {
				r.holdover = data2[limit:]
				data2 = data2[:limit]
			}
			// Append data2 to data
			data = append(data, data2...)
		}
	}

	if !ok {
		return 0, io.EOF
	}

	n = copy(p, data)
	return n, nil
}

// Close closes the reader and the underlying channels.
func (r *XZReader) Close() error {
	// The underying channels close themselves and the memory grabbed by C is freed internally,
	// so we really don't need to do anything here.   It's just good practice to have a close method.
	return nil
}

// NewReader creates a new XZReader that reads from r.   The data represented by r should have been compressed with XZ or LZMA.
func NewReader(r io.Reader) *XZReader {
	return &XZReader{Reader: r, inputchan: make(chan []byte), outputchan: make(chan []byte)}
}
