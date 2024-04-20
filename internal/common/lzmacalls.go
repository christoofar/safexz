// Exposes the internal functions of lzma as a construction of calls requiring streaming
// input and output channels.
package internal

// Connects the input and output channels to the lzma compression function.  When the output channel is closed
// then the compression process is complete.  You must close the input channel to signal the end of the input stream.
func CompressIn(in chan []byte, out chan []byte, strategy int) {

	// The reason for the nested go routines is to isolate the unsafeBuffer
	go func(input chan []byte, output chan []byte, strategy int) {
		go func(receive chan []byte, sender chan []byte, strategy int) {
			compressChanStream(receive, sender, strategy)
		}(input, output, strategy)
	}(in, out, strategy)

}

// Connects the input and output channels to the lzma decompression function.  When the output channel is closed
// then the decompression process is complete.  You must close the input channel to signal the end of the input stream.
func DecompressIn(in chan []byte, out chan []byte) {

	// The reason for the nested go routines is to isolate the unsafeBuffer
	go func(input chan []byte, output chan []byte) {
		go func(receive chan []byte, sender chan []byte) {
			decompressChanStream(receive, sender)
		}(input, output)
	}(in, out)

}
