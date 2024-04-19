// Exposes the internal functions of lzma as a construction of calls requiring streaming
// input and output channels.
package internal

func CompressIn(in chan []byte, out chan []byte, strategy int) {

	// The reason for the nested go routines is to isolate the unsafeBuffer
	go func(input chan []byte, output chan []byte, strategy int) {
		go func(receive chan []byte, sender chan []byte, strategy int) {
			compressChanStream(receive, sender, strategy)
		}(input, output, strategy)
	}(in, out, strategy)

}
