/* trunk-ignore-all(golangci-lint/typecheck) */
// Description: This file contains the functions that connect the input and output channels to the lzma compression and decompression functions.
package internal

import "fmt"

// Connects the input and output channels to the lzma compression function.  When the output channel is closed
// then the compression process is complete.  You must close the input channel to signal the end of the input stream.
func CompressIn(in chan []byte, out chan []byte, strategy int) error {

	errchan := make(chan bool, 1)

	// The reason for the nested go routines is to isolate the unsafeBuffer
	go func(input chan []byte, output chan []byte, strategy int) {
		go func(receive chan []byte, sender chan []byte, strategy int) {
			CompressChanStream(receive, sender, strategy, errchan)
		}(input, output, strategy)
	}(in, out, strategy)

	if <-errchan {
		return fmt.Errorf("liblzma compression failed")
	}

	return nil
}

// Connects the input and output channels to the lzma decompression function.  When the output channel is closed
// then the decompression process is complete.  You must close the input channel to signal the end of the input stream.
func DecompressIn(in chan []byte, out chan []byte) {

	// The reason for the nested go routines is to isolate the unsafeBuffer
	go func(input chan []byte, output chan []byte) {
		go func(receive chan []byte, sender chan []byte) {
			DecompressChanStream(receive, sender)
		}(input, output)
	}(in, out)

}
