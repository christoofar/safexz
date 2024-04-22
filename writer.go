package safexz

import (
	"io"

	internal "github.com/christoofar/safexz/internal"
)

type XZWriter struct {
	io.Writer                      // Support the io.Writer interface
	inputchan  chan []byte         // From /internal, this is where we send the uncompressed data
	outputchan chan []byte         // From /internal, this is where we receive the compressed data
	started    bool                // We need to know if we've started the compressor so we can start it only once.
	done       chan bool           // We need to know when the compressor is done so we can close the output channel.
	Strategy   CompressionStrategy // The compression strategy to use
	funcerr    error               // This is gross, but it's a way to bubble up an lzma crash to topside
}

// Write takes uncompressed data passed in from the underlying source and yields the LZMA2 compressed data into a byte slice.
// You must call Close when you are done writing data to the writer to signal to LZMA there is more data coming (if your stream source is a buffer and not sending EOF)
// otherwise your program will hang.
func (w *XZWriter) Write(p []byte) (n int, err error) {
	// This if block is a run-once context for the goroutines that will hitch to liblzma.  This is because
	// the ABI expects you to call Write multiple times, but we only want to start the compressor once.
	if !w.started {
		w.done = make(chan bool)
		// Start the compessor.  This has to be done in a goroutine here because we will hang execution (internals is written to use gor to stream)
		go func() {
			w.funcerr = internal.CompressIn(w.inputchan, w.outputchan, int(w.Strategy))
		}()

		// Start the output channel.  Whatever is coming out of the compresser is written to the underlying writer.
		// When we're done ranging over the output from lzma, we signal that we're done.
		go func() {
			for data := range w.outputchan {
				if len(data) > 0 {
					w.Writer.Write(data)
				}
			}
			w.done <- true
		}()
		w.started = true // Run-once flag
	}

	// A nil slice is the signal to the compressor that the input stream is done.  This is the only way to signal the end of the input stream.
	// If you don't do this, the compressor will hang waiting for more data.  (Specifically, the close of w.inputchan is the signal to the compressor that the input stream is done.)
	if len(p) == 0 {
		close(w.inputchan)
		w.started = false
	}

	// Send the data to the compressor in 1024-byte blocks.  This is the pattern the ABI expects, no matter
	// what the size of the incoming was.
	for i := 0; i < len(p); i += 1024 {
		end := i + 1024
		if end > len(p) {
			end = len(p)
		}
		w.inputchan <- p[i:end]
	}

	// All we can do here is acknowledge that we got the command to write the data.  Which is what the ABI expects you to say here.
	// The actual compression is happening in a separate go routine.
	//
	// God do I hate the ByteReader/ByteWriter pattern, but it's what the ABI expects.
	return len(p), w.funcerr
}

// Close closes the writer and the underlying channels.  If you do not call Close when you are through writing, LZMA will assume you have more
// data coming and your program will hang.  Note: On big data streams, the LZMA closing process could take quite a while.  If this is bothering you,
// consider using a goroutine to call Close.   In htop or other task monitors you will see the lzma tasks dying off and releasing memory.
func (w *XZWriter) Close() error {
	// We need to hold this close from returning until the output channel is closed.  This is so hokey but it's the crappy ByteReader/ByteWriter
	// pattern at its most luxuriant.

	w.Write([]byte{}) // Send a nil slice to the compressor to signal the end of the input stream.
	<-w.done          // Wait for the compressor to finish
	return nil
}

// NewWriter creates a new XZWriter that writes to w.   The data written to w will be compressed with XZ, yielding an LZMA2 stream.
func NewWriter(w io.Writer, strategy ...CompressionStrategy) *XZWriter {
	use_strategy := CompressionMulti
	if len(strategy) > 0 {
		use_strategy = strategy[0]
	}
	return &XZWriter{Writer: w, inputchan: make(chan []byte), outputchan: make(chan []byte), Strategy: use_strategy}
}
